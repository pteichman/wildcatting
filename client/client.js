
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

function lobby(game) {
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

    return survey;
}

function survey(game) {
    $("#survey").hide()
    $("#report").show()

    return report;
}

function report(game) {
    $("#report").hide()
    $("#drill").show()

    return sell;
}

function drill(game) {
    $("#drill").hide()
    $("#sell").show()

    return sell;
}

function sell(game) {
    $("#sell").hide()
    $("#lobby").show()

    return lobby
}

var state = lobby;

function move() {
    state = state();
}

function prob() {
    $("#prob").hide()
    $("#cost").show()
    return cost;
}

function cost() {
    $("#cost").hide()
    $("#tax").show()
    return tax;
}

function tax() {
    $("#tax").hide()
    $("#oil").show()
    return oil;
}

function oil() {
    $("#oil").hide()
    $("#prob").show()
    return prob;
}

var view = prob;

function tab() {
    view = view();
}

Mousetrap.bind('space', move);
Mousetrap.bind('tab', tab);
