package me.modmuss50.cursemapper

import javafx.application.Application
import javafx.stage.Stage
import javafx.scene.chart.CategoryAxis
import javafx.scene.chart.NumberAxis
import javafx.scene.chart.LineChart

import javafx.scene.chart.XYChart
import javafx.scene.Scene
import javafx.stage.FileChooser
import org.apache.commons.io.FileUtils
import java.nio.charset.Charset


class Graph : Application() {

    override fun start(stage: Stage) {
        stage.title = "Curse data view"
        val xAxis = CategoryAxis()
        val yAxis = NumberAxis()
        xAxis.label = "Date & Time"

        val fileChooser = FileChooser() // Pick the input file
        fileChooser.title = "Open Export File"
        val lines = FileUtils.readLines(fileChooser.showOpenDialog(stage), Charset.defaultCharset())

        yAxis.isForceZeroInRange = false //makes it not start at 0

        val lineChart = LineChart(xAxis, yAxis)
        lineChart.title = "Total Downloads"
        val series = XYChart.Series<String, Number>()
        series.name = "Downloads"

        lines //adds the points
                .map { it.split(",") }
                .forEach { series.data.add(XYChart.Data(it[0].split(" ")[1], it[1].toInt())) }

        val scene = Scene(lineChart, 800.0, 600.0)
        lineChart.data.add(series)
        stage.scene = scene
        stage.show()
    }

    companion object {
        @JvmStatic
        fun start(args: Array<String>) {
            launch(Graph::class.java)
        }
    }
}


fun main(args: Array<String>) {
    Graph.start(args)
}