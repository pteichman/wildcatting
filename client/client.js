
var players = ["bob","joe","jbz","peter","snorp"];

var probColor = d3.scale.quantize().domain([1, 100]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var costColor = d3.scale.quantize().domain([1, 25]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var taxColor = d3.scale.quantize().domain([100, 550]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var oilColor = d3.scale.quantize().domain([1, 9]).range(["#4d4d4d",,"#878787","#bababa","#e0e0e0","#ffffff","#fddbc7","#f4a582","#d6604d","#b2182b"]);

var game = {site: 0, data: {}};

var fsm = StateMachine.create({

    initial: 'lobby',

    events: [
        { name: 'play', from: 'lobby', to: 'survey' },
        { name: 'survey', from: 'survey', to: 'report' },
        { name: 'yes', from: 'report', to: 'drill' },
        { name: 'no', from: 'report', to: 'sell' },
        { name: 'done', from: 'drill', to: 'sell' },
        { name: 'done', from: 'sell', to: 'score' },
        { name: 'done', from: 'score', to: 'survey' },
    ],

    callbacks: {

        onenterlobby: lobby,
        onentersurvey: survey,
        onenterreport: report,
        onenterdrill: drill,
        onentersell: sell,
        onenterscore: score,

        onleavelobby: function() {
            $("#lobby").hide();
            Mousetrap.reset();
        },
        onleavesurvey: function() {
            $("#survey").hide();
            Mousetrap.reset();
        },
        onleavereport: function() {
            $("#report").hide();
            Mousetrap.reset();
        },
        onleavedrill: function() {
            $("#drill").hide();
            Mousetrap.reset();
        },
        onleavesell: function() {
            $("#sell").hide();
            Mousetrap.reset();
        },
        onleavescore: function() {
            $("#score").hide();
            Mousetrap.reset();
        },
    }
});

function lobby() {
    $("#lobby").show();

    d3.select("#lobby-players")
        .selectAll("p")
        .data(players)
        .enter()
        .append("p")
        .text(function(d) { return d; });

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        $.post("/game/0/player/0/", '{"done":true}', function(data) {
            game.data = data;
            fsm.play();
        });
    });
}

function survey() {
    $("#survey").show();

    d3.select("#prob")
        .selectAll("rect")
        .data(game.data.prob)
        .enter()
        .append("rect")
        .attr("class", function (d, i) { return "site-" + i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", probColor);

    d3.select("#cost")
        .selectAll("rect")
        .data(game.data.cost)
        .enter()
        .append("rect")
        .attr("class", function (d, i) { return "site-" + i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", costColor);

    d3.select("#tax")
        .selectAll("rect")
        .data(game.data.tax)
        .enter()
        .append("rect")
        .attr("class", function (d, i) { return "site-" + i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", taxColor);

    d3.select("#oil")
        .selectAll("rect")
        .data(game.data.oil)
        .enter()
        .append("rect")
        .attr("class", function (d, i) { return "site-" + i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", function (d) { return d == 0 ? 'black' : oilColor(d); });

    d3.select("#fact").text(game.data.fact);
    d3.select("#week").text("Week " + game.data.week)

    $(".site-0").toggleClass("cursor");

    var views = ["#prob", "#cost", "#tax", "#oil"];
    var cur = 0;
    function view(delta) {
        $(views[cur]).hide();
        cur = mod(cur+delta, views.length);
        $(views[cur]).show();
    }

    function cursor(dy, dx) {
        $(".site-" + game.site).toggleClass("cursor");

        var y = mod(Math.floor(game.site/80)+dy, 24);
        var x = mod(mod(game.site, 80)+dx, 80);
        game.site = y*80 + x;

        $(".site-" + game.site).toggleClass("cursor");
    }

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        fsm.survey();
    });
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
    });
    Mousetrap.bind('up', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(-1, 0);
    });
}

function report() {
    $("#report").show()

    // TODO define an Object to serialize you weirdo
    $.post("/game/0/player/0/", '{"site":' + game.site + '}', function(data) {
        d3.select("#report-site").text("X="+mod(game.site, 80)+"\tY="+Math.floor(game.site/80));
        d3.select("#report-prob").text(data.prob + "%");
        d3.select("#report-cost").text("$\t" + data.cost);
        d3.select("#report-tax").text("$\t" + data.tax);
    }).done(function() {
        Mousetrap.bind('y', function(e) {
            e.preventDefault ? e.preventDefault() : (e.returnValue = false);
            fsm.yes()
        });
        Mousetrap.bind('n', function(e) {
            e.preventDefault ? e.preventDefault() : (e.returnValue = false);
            fsm.no()
        });
    });
}

function drill() {
    $("#drill").show()

    function advance() {
        $.post("/game/0/player/0/", '{}', function(data) {
            if (data.oil) {
                fsm.done()
            } else if (data.Depth == 9) {
                fsm.done()
            }
        });
    }

    advance();

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        advance();
    });

    Mousetrap.bind('q', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        $.post("/game/0/player/0/", '{"done":true}', function(data) {
            fsm.done();
        });
    });
}

function sell() {
    $("#sell").show()
    Mousetrap.bind('q', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        $.post("/game/0/player/0/", '{"done":true}', function(data) {
            fsm.done()
        });
    });
}

function score() {
    $("#score").show()
    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        fsm.done()
    });
}

// % operator in javascript is remainder and isn't helpful for wrapping negatives
function mod(a, n) {
    return a - (n * Math.floor(a/n));
}
