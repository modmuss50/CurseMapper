package me.modmuss50.cursemapper

import org.apache.commons.io.FileUtils
import java.io.File
import java.nio.charset.Charset
import java.text.SimpleDateFormat
import java.util.*
import java.util.TimerTask

fun main(args : Array<String>) {
    val projects = File("projects")
    if(!projects.exists()){
        projects.mkdir()
    }
    val timer = Timer()
    timer.schedule(object : TimerTask() {
        override fun run() {
            println("Loading data")
            for(user in FileUtils.readLines(File("users.txt"), Charset.defaultCharset())){
                println("loading data for ${user}")
                val curseData = CurseUtil(user)
                curseData.load()
                FileUtils.writeStringToFile(File("${user}_export.csv"), "${SimpleDateFormat("dd/MM/yyyy HH:mm:ss").format(Date())},${curseData.totalDownloads}${System.getProperty("line.separator")}", Charset.defaultCharset(), true)
                var userFolder = File(projects, user)
                if(!userFolder.exists()){
                    userFolder.mkdir()
                }
                for (set in curseData.downloads!!.entries) {
                    FileUtils.writeStringToFile(File(userFolder, "${set.key.substring(set.key.indexOf("[") + 1,set.key.indexOf("]")).replace(" ", "_")}_export.csv"), "${SimpleDateFormat("dd/MM/yyyy HH:mm:ss").format(Date())},${set.value}${System.getProperty("line.separator")}", Charset.defaultCharset(), true)
                }
            }
            println("Done")
        }
    }, 0, 1000 * 60 * 15) // 1000 milliseconds in a second, 60 seconds in a min, and 15 mins in 15 mins
}