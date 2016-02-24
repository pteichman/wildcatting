
var players = ["bob","joe","jbz","peter","snorp"];

var probColor = d3.scale.quantize().domain([1, 100]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var costColor = d3.scale.quantize().domain([1, 25]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var taxColor = d3.scale.quantize().domain([100, 550]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var oilColor = d3.scale.quantize().domain([1, 9]).range(["#4d4d4d",,"#878787","#bababa","#e0e0e0","#ffffff","#fddbc7","#f4a582","#d6604d","#b2182b"]);

function play() {
    d3.select("#lobby-players")
        .selectAll("p")
        .data(players)
        .enter()
        .append("p")
        .text(function(d) { return d; });

    return lobby
}

function lobby() {
    $("#lobby").hide()

    $.post("/game/0/player/0/", '{"done":true}', function(data) {
        d3.select("#prob")
            .selectAll("rect")
            .data(data.prob)
            .enter()
            .append("rect")
            .attr("class", function (d, i) { return "site-" + i; })
            .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
            .attr("x", function (d, i) { return i%80 * 12 ; })
            .style("fill", probColor);

        d3.select("#cost")
            .selectAll("rect")
            .data(data.cost)
            .enter()
            .append("rect")
            .attr("class", function (d, i) { return "site-" + i; })
            .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
            .attr("x", function (d, i) { return i%80 * 12 ; })
            .style("fill", costColor);

        d3.select("#tax")
            .selectAll("rect")
            .data(data.prob)
            .enter()
            .append("rect")
            .attr("class", function (d, i) { return "site-" + i; })
            .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
            .attr("x", function (d, i) { return i%80 * 12 ; })
            .style("fill", taxColor);

        d3.select("#oil")
            .selectAll("rect")
            .data(data.prob)
            .enter()
            .append("rect")
            .attr("class", function (d, i) { return "site-" + i; })
            .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
            .attr("x", function (d, i) { return i%80 * 12 ; })
            .style("fill", function (d) { return d == 0 ? 'black' : oilColor(d); });

        d3.select("#fact").text(data.fact);
        d3.select("#week").text("Week " + data.week)
    }).done(function () {
        if ($(".cursor").length == 0) {
            $(".site-0").toggleClass("cursor");
        }
        $("#survey").show()
    });
    return survey;
}

function survey() {
    $("#survey").hide()
    $("#report").show()

    // figure out the cursor site
    // ...

    $.post("/game/0/player/0/", '{"site":1}', function(data) {
        d3.select("#report-site").text(data.site);
        d3.select("#report-prob").text(data.prob + "%");
        d3.select("#report-cost").text("$\t" + data.cost);
        d3.select("#report-tax").text("$\t" + data.tax);
    });

    return report;
}

function report() {
    $("#report").hide()
    $("#drill").show()

    return drill;
}

function drill() {
    $("#drill").hide()
    $("#sell").show()

    return sell;
}

function sell() {
    $("#sell").hide()
    $("#lobby").show()

    return lobby;
}

// % operator in javascript is remainder and isn't helpful for wrapping negatives
function mod(a, n) {
    return a - (n * Math.floor(a/n));
}

function playState() {
    var state = play;

    return function() {
        state = state();
    }
}

function viewState() {
    var views = ["#prob", "#cost", "#tax", "#oil"];
    var cur = 0;

    return function(delta) {
        $(views[cur]).hide();
        cur = mod(cur+delta, views.length);
        $(views[cur]).show();
    }
}

function cursorState() {
    var y = 0;
    var x = 0;

    return function(dy, dx) {
        $(".site-" + (y * 80 + x)).toggleClass("cursor");

        y = mod(y+dy, 24)
        x = mod(x+dx, 80)

        $(".site-" + (y * 80 + x)).toggleClass("cursor");
        return site;
    }
}

var play = playState();
var view = viewState();
var cursor = cursorState();

Mousetrap.bind('space', play);
Mousetrap.bind('tab', function(e) {
    e.preventDefault ? e.preventDefault() : (e.returnValue = false);
    view(1);
});
Mousetrap.bind('shift+tab', function(e) {
    e.preventDefault ? e.preventDefault() : (e.returnValue = false);
    view(-1);
});
Mousetrap.bind('left', function(e) {
    e.preventDefault ? e.preventDefault() : (e.returnValue = false);
    cursor(0, -1);
});
Mousetrap.bind('down', function(e) {
    e.preventDefault ? e.preventDefault() : (e.returnValue = false);
    cursor(1, 0);
});
Mousetrap.bind('right', function(e) {
    e.preventDefault ? e.preventDefault() : (e.returnValue = false);
    cursor(0, 1);
})
Mousetrap.bind('up', function(e) {
    e.preventDefault ? e.preventDefault() : (e.returnValue = false);
    cursor(-1, 0);
});

// set everything in motion
play();
