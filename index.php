<!DOCTYPE html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="HandheldFriendly" content="True">
    <meta name="MobileOptimized" content="320">

    <title>Curse Data</title>
    <meta name="description" content="Curse Data">

    <script src="https://code.jquery.com/jquery-3.1.1.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery-csv/0.71/jquery.csv-0.71.min.js"></script>
    <script type="text/javascript" src="https://www.google.com/jsapi"></script>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>

    <link rel="stylesheet"
          href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/4.0.2/bootstrap-material-design.css"/>
    <script
        src="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-material-design/4.0.2/bootstrap-material-design.umd.min.js"></script>

    <link rel="stylesheet" href="css/style.css"/>
    <script type="text/javascript" src="js/chart.js"></script>

</head>

<body>
<div class="container">

<?php if (!isset($_GET["user"])) : ?>




    <div class="panel panel-default">
        <div class="panel-heading">Curse tracker</div>
        <div id="infoPanel" class="panel-body">
            <p>This is a website used to track downloads on curse accounts. Only certain user accounts can have their
                downloads tracked every 15 mins and displayed on a graph.
                This is because the data is scraped from the website and if we make too many requests in a short time we
                will be blocked from doing this.</p>

            <p>Feel free to contact moduss50 if you feel like you should be added, but expect to be turned out if I dont
                know you. Sorry.</p>

            <p>But, wait! You can still lookup your current total downloads on curseforge! You will not be added to the
                tracker, but will be able to see your total downloads on your curse account.</p>

            <form>
                <div class="form-group" id="infoForm">

                    <div class="col-sm-10">
                        <input  class="form-control" id="i5ps" placeholder="Curse username">
                    </div>
                    <button id="lookupButton" type="button" class="btn btn-primary">Go</button>
                </div>
            </form>
        </div>
        <div id="loading" class="panel-body">
            <h2>Coming soon.</h2>
            <p>I've not coded it yet, so hold on. :P</p>

        </div>
    </div>

    <?php endif; ?>



    <?php

    $dirname = ".";
    if (isset($_GET["user"])) {
        $dirname = "./projects/" . $_GET["user"];

        $userFile = "./" . $_GET["user"] . "_export.csv";
        if (file_exists($userFile)) {
            echo '</P><a href="./" class="btn btn-default" role="button">Home Page</a><hr>';
            echo '<p><b>' . $_GET["user"] . '</b> has a total download count of: ' . getTotalUserDownloads($_GET["user"]) . '</p>';
        } else {
            echo '<p>User not found, ask <a href="https://twitter.com/modmuss50">modmuss50</a> to see if its possible to be tracked.</p><p>Please note there are limited places due to the worry of being rate limited by curse.</p>';
            echo '';
            echo '<p><a href="./">Go Home</a></p>';
        }
    }

    if ($handle = opendir($dirname)) {
        while (false !== ($file = readdir($handle))) {
            if ($file != "." && $file != ".." && strtolower(substr($file, strrpos($file, '.') + 1)) == 'csv') {
                $name = $file;

                echo '<div class="panel panel-default"><div class="panel-heading"><h3 class="panel-title pull-left"></h3>';


                if (isset($_GET["user"])) {
                    $name = "./projects/" . $_GET["user"] . DIRECTORY_SEPARATOR . $file;
                    $arr = explode("_", $file, 2);
                    $uname = $arr[0];
                    echo '<p><b>' . $uname . '</b> has a total download count of: ' . getTotalProjectDownloads($_GET["user"], $file) . '</p></h3>';
                } else {
                    $arr = explode("_", $file, 2);
                    $uname = $arr[0];
                    echo '<a href="./?user=' . ${'uname'} . '">View all of <b>' . ${'uname'} . 's</b> projects</a></h3>';
                    echo '<div class="panel-title pull-right">Total downloads = ' . getTotalUserDownloads(${'uname'}) . '</div>';
                }

                echo '</div><div class="panel-body">';

                echo '<div id="' . $name . '"></div></div></div>';


                echo '<script type="text/javascript">';
                echo "google.setOnLoadCallback(drawChart('${name}'));";
                echo '</script>';

                echo '<script type="text/javascript">';
                echo "$(window).resize(function(){drawChart('${name}');});";
                echo '</script>';


            }
        }
        closedir($handle);
    }

    function endswith($string, $test)
    {
        $strlen = strlen($string);
        $testlen = strlen($test);
        if ($testlen > $strlen) return false;
        return substr_compare($string, $test, $strlen - $testlen, $testlen) === 0;
    }

    function getTotalUserDownloads($user)
    {
        $file = "./${'user'}_export.csv";
        $data = file($file);
        $line = $data[count($data) - 1];

        return number_format((int)substr($line, strpos($line, ",") + 1));
    }

    function getTotalProjectDownloads($user, $name)
    {
        $file = "./projects/${'user'}/${'name'}";
        $data = file($file);
        $line = $data[count($data) - 1];

        return number_format((int)substr($line, strpos($line, ",") + 1));
    }

    ?>


</div>
</body>


<footer class="footer">
    <div class="container">
        <p class="text-muted">
            <a href="https://discord.gg/0tCDWb77cvetwm0e"><img border="0" alt="Discord"
                                                               src="https://img.shields.io/badge/Discord-TeamReborn-738bd7.svg"></a>

            <a href="https://github.com/modmuss50/CurseMapper"><img border="0" alt="Discord"
                                                                    src="https://img.shields.io/github/stars/modmuss50/CurseMapper.svg"></a>

            <a href="https://twitter.com/modmuss50"><img border="0" alt="Discord"
                                                         src="https://img.shields.io/twitter/follow/modmuss50.svg?style=social&label=Follow"></a>
        </p>
    </div>
</footer>


<script src="js/script.js"></script>

