
var players = ["bob","joe","jbz","peter","snorp"];

d3.select("#lobby-players")
    .selectAll("p")
    .data(players)
    .enter()
    .append("p")
    .text(function(d) { return d; });

var probColor = d3.scale.quantize().domain([1, 100]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var costColor = d3.scale.quantize().domain([1, 25]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var taxColor = d3.scale.quantize().domain([100, 550]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var oilColor = d3.scale.quantize().domain([1, 9]).range(["#4d4d4d",,"#878787","#bababa","#e0e0e0","#ffffff","#fddbc7","#f4a582","#d6604d","#b2182b"]);

function lobbyState(game) {
    $("#lobby").hide()

    $.post("/game/0/player/0/", '{"done":true}', function(data) {
        d3.select("#prob")
           .append("svg")
           .selectAll("rect")
           .data(data.prob)
           .enter()
           .append("rect")
           .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
           .attr("x", function (d, i) { return i%80 * 12 ; })
           .style("fill", probColor);

       d3.select("#cost")
          .append("svg")
          .selectAll("rect")
          .data(data.cost)
          .enter()
          .append("rect")
          .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
          .attr("x", function (d, i) { return i%80 * 12 ; })
          .style("fill", costColor);

      d3.select("#tax")
         .append("svg")
         .selectAll("rect")
         .data(data.prob)
         .enter()
         .append("rect")
         .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
         .attr("x", function (d, i) { return i%80 * 12 ; })
         .style("fill", taxColor);

     d3.select("#oil")
        .append("svg")
        .selectAll("rect")
        .data(data.prob)
        .enter()
        .append("rect")
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", function (d) { return d == 0 ? 'black' : oilColor(d); });

        d3.select("#fact").text(data.fact);
        d3.select("#week").text("Week " + data.week)
    });

    $("#survey").show()

    return surveyState;
}

function surveyState(game) {
    $("#survey").hide()
    $("#report").show()

    return reportState;
}

function reportState(game) {
    $("#report").hide()
    $("#drill").show()

    return sellState;
}

function drillState(game) {
    $("#drill").hide()
    $("#sell").show()

    return sellState;
}

function sellState(game) {
    $("#sell").hide()
    $("#lobby").show()

    return lobbyState;
}

var curState = lobbyState;

function move() {
    curState = curState();
}

// % operator in javascript is remainder and isn't helpful for wrapping negatives
function mod(a, n) {
    return a - (n * Math.floor(a/n));
}

function viewState() {
    var views = ["#prob", "#cost", "#tax", "#oil"];
    var cur = 0;

    return function(delta) {
        $(views[cur]).hide();
        cur = mod((cur + delta), views.length);
        $(views[cur]).show();
    }
}

var view = viewState();

Mousetrap.bind('space', move);
Mousetrap.bind('tab', function(e) {
    event.preventDefault ? event.preventDefault() : (event.returnValue = false);
    return view(1)
    });
Mousetrap.bind('shift+tab', function(e) {
    event.preventDefault ? event.preventDefault() : (event.returnValue = false);
    return view(-1);
    });
