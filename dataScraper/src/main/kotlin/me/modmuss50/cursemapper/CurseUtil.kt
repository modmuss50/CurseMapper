package me.modmuss50.cursemapper

import java.io.BufferedReader
import java.io.FileNotFoundException
import java.io.InputStreamReader
import java.net.MalformedURLException
import java.net.URL
import java.util.*

/**
 * quick port to koltin: modmuss50
 */


// Based on code by Jared found here: https://github.com/MinecraftModDevelopment/MMDBot/blob/master/src/main/java/com/mcmoddev/bot/util/CurseData.java
class CurseUtil(val username: String) {

    var avatar: String? = null
        private set

    var profile: String = ""

    var projectURLs: MutableList<String>? = null

    var downloads: MutableMap<String, Long>? = null

    var totalDownloads: Long = 0

    var foundUser = true

    var foundProject = false

    fun load() {
        profile = PROFILE_URL + username
        this.projectURLs = ArrayList<String>()
        this.downloads = HashMap<String, Long>()

        val baseURL = String.format(PROJECTS_PAGE_URL, username)

        var page = 1
        var foundDuplicatePage = false
        var foundAvatar = false

        while (!foundDuplicatePage) {

            try {
                URL(baseURL + page).openStream().use { stream ->

                    val reader = BufferedReader(InputStreamReader(stream))

                    var line: String? = reader.readLine()
                    while (line != null) {
                        // Find avatar
                        if (!foundAvatar && line.contains("<div class=\"avatar avatar-100 user user-role-curse-premium\">")) {

                            reader.readLine()
                            val imageLine = reader.readLine()

                            if (imageLine.contains("<img src=")) {

                                this.avatar = imageLine.split("\"".toRegex()).dropLastWhile { it.isEmpty() }.toTypedArray()[1]
                                foundAvatar = true
                            }
                        } else if (line.contains("a href=\"/projects/")) {

                            val projectURL = "https://minecraft.curseforge.com" + line.split("\"".toRegex()).dropLastWhile { it.isEmpty() }.toTypedArray()[1].split("\"".toRegex()).dropLastWhile { it.isEmpty() }.toTypedArray()[0]

                            if (!this.projectURLs!!.contains(projectURL)) {

                                this.projectURLs!!.add(projectURL)
                                this.foundProject = true
                            } else {

                                foundDuplicatePage = true
                                break
                            }
                        }// Find projects
                        line = reader.readLine()
                    }

                    if (page == 1 && !this.foundProject)
                        return
                }
            } catch (e: Exception) {

                if (e is java.io.IOException && e.message!!.contains("HTTP response code: 400") || e is MalformedURLException || e is FileNotFoundException && (e.message!!.contains("https://minecraft.curseforge.com/not-found?404") || e.message!!.contains("https://minecraft.curseforge.com/members/"))) {

                    this.foundUser = false
                    break
                } else
                    e.printStackTrace()
            }

            page++
        }

        if (!this.foundUser || !this.foundProject)
            return

        // downloads
        for (projectUrl in this.projectURLs!!)
            try {
                URL(projectUrl).openStream().use { stream ->

                    val reader = BufferedReader(InputStreamReader(stream))
                    var foundDownloads = false

                    var line: String? = reader.readLine()
                    while (line != null) {
                        if (foundDownloads) {

                            val projectDownloads = java.lang.Long.parseLong(line!!.split(">".toRegex()).dropLastWhile { it.isEmpty() }.toTypedArray()[1].split("<".toRegex()).dropLastWhile { it.isEmpty() }.toTypedArray()[0].replace(",".toRegex(), ""))
                            this.totalDownloads += projectDownloads
                            val name = projectUrl.replace("https://minecraft.curseforge.com/projects/", "").replace("-".toRegex(), " ")

                            this.downloads!!.put("[" + name + "](" + projectUrl.replace(" ".toRegex(), "-") + ")", projectDownloads)
                            break
                        } else if (line!!.contains("Total Downloads"))
                            foundDownloads = true
                        line = reader.readLine()
                    }
                }
            } catch (e: Exception) {

                e.printStackTrace()
            }
    }


    fun hasProjects(): Boolean {

        return this.foundProject && !this.projectURLs!!.isEmpty() && !this.downloads!!.isEmpty()
    }

    val projectCount: Int
        get() = this.projectURLs!!.size

    companion object {

        val PROFILE_URL = "https://minecraft.curseforge.com/members/"

        val PROJECTS_PAGE_URL = "https://minecraft.curseforge.com/members/%s/projects?page="

    }

}
