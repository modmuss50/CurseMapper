<!DOCTYPE html>
<head>
 <meta charset="utf-8">
 <meta http-equiv="X-UA-Compatible" content="IE=edge">
 <meta name="viewport" content="width=device-width, initial-scale=1">
 <meta name="HandheldFriendly" content="True">
 <meta name="MobileOptimized" content="320">

 <title>Curse Data</title>
 <meta name="description" content="Curse Data">

 <script src="http://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
 <script src="https://code.jquery.com/jquery-3.1.1.min.js"></script>
 <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery-csv/0.71/jquery.csv-0.71.min.js"></script>
 <script type="text/javascript" src="http://www.google.com/jsapi"></script>

<script type="text/javascript">
  google.load('visualization', '1', { packages: ['corechart', 'controls'] });
</script>

<script type="text/javascript">

    function drawChart(username) {
$.ajaxSetup({ cache: false });
        
        $.get( username, function(csvString) {
      
        var arrayData = $.csv.toArrays("Date,Downloads\n" + csvString, {onParseValue: $.csv.hooks.castToScalar});
        var data = new google.visualization.arrayToDataTable(arrayData);
        var view = new google.visualization.DataView(data);
        view.setColumns([0,1]);
        var graphtitle = /[^/]*$/.exec(username)[0].replace("_export.csv", "");
        var options = {
            title: graphtitle,
            hAxis: {title: data.getColumnLabel(0), minValue: data.getColumnRange(0).min, maxValue: data.getColumnRange(0).max},
            vAxis: {title: data.getColumnLabel(1), minValue: data.getColumnRange(1).min, maxValue: data.getColumnRange(1).max},
            legend: 'none'
        };
        var chart = new google.visualization.LineChart(document.getElementById(username));
        chart.draw(view, options);
    });
}
</script>

<?php

            $dirname = ".";
            if(isset($_GET["user"])){
                $dirname = "./projects/" . $_GET["user"];
                echo '<p>'. $_GET["user"] . ' has a total download count of: ' . getTotalUserDownloads($_GET["user"]) . '</p>';
            }

            if ($handle = opendir($dirname)) {
                while (false !== ($file = readdir($handle)))
                {
                    if ($file != "." && $file != ".." && strtolower(substr($file, strrpos($file, '.') + 1)) == 'csv')
                    {
                        $name = $file;
                        if(isset($_GET["user"])){
                            $name = "./projects/" .$_GET["user"] . DIRECTORY_SEPARATOR . $file;
                        } 
                        else {
                            $arr = explode("_", $file, 2);
                            $uname = $arr[0];
                            echo  '<a href="./?user=' .${'uname'} .'">View all of ' .${'uname'} .'s projects</a>';
                            echo '<p>'. ${'uname'} . ' has a total download count of: ' . getTotalUserDownloads(${'uname'}) . '</p>';
                        }
                        echo  '<div id="' . $name . '"></div>';
                        echo  '<script type="text/javascript">';
                        echo  "google.setOnLoadCallback(drawChart('${name}'));";
                       echo'</script>';
                       
                    }
                }
                closedir($handle);
            }

            function endswith($string, $test) {
                $strlen = strlen($string);
                $testlen = strlen($test);
                if ($testlen > $strlen) return false;
                return substr_compare($string, $test, $strlen - $testlen, $testlen) === 0;
            }

            function getTotalUserDownloads($user){
                $file = "./${'user'}_export.csv";
                $data = file($file);
                $line = $data[count($data)-1];
                
                return number_format((int)substr($line, strpos($line, ",") + 1));
            }
            
        ?>
