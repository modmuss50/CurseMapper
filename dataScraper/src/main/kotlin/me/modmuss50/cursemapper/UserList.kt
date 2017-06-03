package me.modmuss50.cursemapper

import com.google.common.collect.Ordering
import com.google.gson.GsonBuilder
import com.sun.deploy.util.ArrayUtil
import me.modmuss50.ArraySorter
import org.apache.commons.io.FileUtils
import org.apache.commons.io.IOUtils
import java.io.BufferedReader
import java.io.File
import java.io.InputStreamReader
import java.net.URL
import java.util.*
import java.util.regex.Pattern
import java.io.StringWriter
import java.nio.charset.Charset
import java.util.ArrayList

class UserList {
    var projects = ArrayList<String>()
    var users = TreeMap<String, Int>()

    fun load() {
        val pages = 100
        for (i in 1..pages) {
            println(i)
            URL(PROJECTS_URL + 1).openStream().use { stream ->
                val reader = BufferedReader(InputStreamReader(stream))
                var line: String? = reader.readLine()
                while (line != null) {
                    if(line.contains("<a href=") && line.contains("/projects/")){
                        val r = Pattern.compile("(?<=\")(.*?)(?=\")")
                        val m = r.matcher(line)
                        m.find()
                        projects.add(m.group(0).split("/")[2])
                        loadProject(m.group(0).split("/")[2])
                    }
                    line = reader.readLine()
                }
            }
        }
    }

    fun loadProject(name : String){
        URL(API_URL + name).openStream().use { stream ->
            val writer = StringWriter()
            IOUtils.copy(stream, writer, Charset.defaultCharset())
            val json = writer.toString()
            val data = GSON.fromJson(json, JsonFomat::class.java)
            for(entry in data.authors){
                if(users.containsKey(entry.key)){
                    users.set(entry.key, (data.download + users[entry.key] as Int))
                } else {
                    users.put(entry.key, data.download)
                }
            }

        }
    }

    companion object {

        val PROJECTS_URL = "https://minecraft.curseforge.com/mc-mods?filter-sort=downloads&page="

        val API_URL = "http://mcmoddev.com/curse/cursereader.php?c="

        val GSON = GsonBuilder().create()
    }

}

class JsonFomat {
    var download : Int = 0
    var authors = HashMap<String, String>()
}


fun main(args: Array<String>){
    val ul = UserList()
    ul.load()
    //ul.loadProject("mantle")
    print(ArraySorter.entriesSortedByValues(ul.users).reversed())

    var i = 1
    for(data in ArraySorter.entriesSortedByValues(ul.users).reversed()){
        FileUtils.writeStringToFile(File("list.txt"), i.toString() +  ":" + data.key + ":" + data.value + System.getProperty("line.separator"), true)
        i++
    }
}

