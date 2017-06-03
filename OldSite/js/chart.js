google.load('visualization', '1', { packages: ['corechart', 'controls'] });

function drawChart(username) {

    $.ajaxSetup({cache: false});

    $.get(username, function (csvString) {

        var arrayData = $.csv.toArrays("Date,Downloads\n" + csvString, {onParseValue: $.csv.hooks.castToScalar});
        var data = new google.visualization.arrayToDataTable(arrayData);
        var view = new google.visualization.DataView(data);
        view.setColumns([0, 1]);
        var graphtitle = /[^/]*$/.exec(username)[0].replace("_export.csv", "");
        var options = {
            title: graphtitle,
            hAxis: {
                title: data.getColumnLabel(0),
                minValue: data.getColumnRange(0).min,
                maxValue: data.getColumnRange(0).max
            },
            vAxis: {
                title: data.getColumnLabel(1),
                minValue: data.getColumnRange(1).min,
                maxValue: data.getColumnRange(1).max
            },
            legend: 'none'
        };
        var chart = new google.visualization.LineChart(document.getElementById(username));
        chart.draw(view, options);
    });
}