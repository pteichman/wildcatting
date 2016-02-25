
var players = ["bob","joe","jbz","peter","snorp"];

var probColor = d3.scale.quantize().domain([1, 100]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var costColor = d3.scale.quantize().domain([10, 250]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var taxColor = d3.scale.quantize().domain([100, 550]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var oilColor = d3.scale.quantize().domain([1, 9]).range(["#4d4d4d",,"#878787","#bababa","#e0e0e0","#ffffff","#fddbc7","#f4a582","#d6604d","#b2182b"]);

var state = {};

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
            d3.select("#lobby").style("display", "none");
            Mousetrap.reset();
        },
        onleavesurvey: function() {
            d3.select("#survey").style("display", "none");
            d3.select("#field svg").selectAll("*").remove();
            Mousetrap.reset();
        },
        onleavereport: function() {
            d3.select("#report").style("display", "none");
            Mousetrap.reset();
        },
        onleavedrill: function() {
            d3.select("#drill").style("display", "none");
            Mousetrap.reset();
        },
        onleavesell: function() {
            d3.select("#sell").style("display", "none");
            d3.select("#sell-table tbody").html("");
            Mousetrap.reset();
        },
        onleavescore: function() {
            d3.select("#score").style("display", "none");
            Mousetrap.reset();
        },
    }
});

function lobby() {
    d3.select("#lobby").style("display", "block");

    d3.select("#lobby-players")
        .selectAll("p")
        .data(players)
        .enter()
        .append("p")
        .text(function(d) { return d; });

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json("/game/0/player/0/")
            .on("load", function(data) {
                state = data;
                fsm.play();
            })
            .on("error", function(error) {
                alert(error);
            })
            .post(JSON.stringify({done: true}));
    });
}

function survey() {
    d3.select("#survey").style("display", "block");

    d3.select("#prob")
        .selectAll("rect")
        .data(state.prob)
        .enter()
        .append("rect")
        .attr("data-site", function (d, i) { return i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", probColor);

    d3.select("#cost")
        .selectAll("rect")
        .data(state.cost)
        .enter()
        .append("rect")
        .attr("data-site", function (d, i) { return i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", costColor);

    d3.select("#tax")
        .selectAll("rect")
        .data(state.tax)
        .enter()
        .append("rect")
        .attr("data-site", function (d, i) { return i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", taxColor);

    d3.select("#oil")
        .selectAll("rect")
        .data(state.oil)
        .enter()
        .append("rect")
        .attr("data-site", function (d, i) { return i; })
        .attr("y", function (d, i) { return Math.floor(i/80) * 18; })
        .attr("x", function (d, i) { return i%80 * 12 ; })
        .style("fill", function (d) { return d == 0 ? 'black' : oilColor(d); });

    d3.select("#fact").text(state.fact);
    d3.select("#week").text("Week " + state.week)

    d3.selectAll("rect[data-site='0'").attr("class", "cursor");

    var views = ["#prob", "#cost", "#tax", "#oil"];
    var cur = 0;
    function view(delta) {
        d3.select(views[cur]).style("display", "none");
        cur = mod(cur+delta, views.length);
        d3.select(views[cur]).style("display", "block");
    }

    var site = 0;
    function cursor(dy, dx) {
        d3.selectAll("rect[data-site='"+site+"']").attr("class", "");

        var y = mod(Math.floor(site/80)+dy, 24);
        var x = mod(mod(site, 80)+dx, 80);
        site = y*80 + x;

        d3.selectAll("rect[data-site='"+site+"']").attr("class", "cursor");
    }

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);

        d3.json("/game/0/player/0/")
            .on("load", function(data) {
                state = data;
                state.site = site; // bleh
                fsm.survey();
            })
            .on("error", function(error) {
                alert(error);
             })
            .post(JSON.stringify({site: site}));
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
    Mousetrap.bind('shift+left', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(0, -3);
    });
    Mousetrap.bind('down', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(1, 0);
    });
    Mousetrap.bind('shift+down', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(3, 0);
    });
    Mousetrap.bind('right', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(0, 1);
    });
    Mousetrap.bind('shift+right', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(0, 3);
    });
    Mousetrap.bind('up', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(-1, 0);
    });
    Mousetrap.bind('shift+up', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        cursor(-3, 0);
    });
}

function report() {
    d3.select("#report").style("display", "block");
    d3.select("#report-site").text("X="+mod(state.site, 80)+"\tY="+Math.floor(state.site/80));
    d3.select("#report-prob").text(state.prob[state.site] + "%");
    d3.select("#report-cost").text("$\t" + state.cost[state.site]);
    d3.select("#report-tax").text("$\t" + state.tax[state.site]);

    Mousetrap.bind('y', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        fsm.yes();
    });
    Mousetrap.bind('n', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        fsm.no();
    });
}

function drill() {
    d3.select("#drill").style("display", "block");

    advance();

    var site = state.site;
    var bit = 0;
    function advance() {
        d3.json("/game/0/player/0/")
            .on("load", function(data) {
                bit++;
                state = data;

                // this is awkward. deeds should be a site->deed map
                for (var i=0; i < state.deeds.length; i++) {
                    if (state.deeds[i].site != site) {
                        continue;
                    }
                    // split these and add DRY HOLE and gusher animations
                    if (state.deeds[i].oil || bit == 9) {
                        fsm.done();
                    }
                }
            })
            .on("error", function(error) {
                alert(error);
             })
            .post(JSON.stringify({}));
    }

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        advance();
    });

    Mousetrap.bind('q', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json("/game/0/player/0/")
            .on("load", function(data) {
                state = data;
                fsm.done();
            })
            .on("error", function(error) { alert(error); })
            .post(JSON.stringify({done: true}));
    });
}

function sell() {
    d3.select("#sell").style("display", "block");

    function siteData(d) {
        var depth = d.oil > 0 ? (d.oil * 100): "";
        var x = mod(d.site, 80);
        var y = Math.floor(d.site/80);
        var tax = state.tax[d.site];
        var pnl = d.income - d.cost;
        var data = [x, y, depth, "$", d.cost, "$", tax, "$", d.income, "$", pnl]
        console.log(data);
        return data;
    }

    d3.select("#sell-table tbody")
        .selectAll("tr")
        .data(state.deeds)
        .enter()
        .append("tr")
        .selectAll("td")
        .data(siteData)
        .enter()
        .append("td")
        .text(function(d) { return d; });

    Mousetrap.bind('q', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json("/game/0/player/0/")
            .on("load", function(data) {
                state = data;
                fsm.done();
            })
            .on("error", function(error) { alert(error); })
            .post(JSON.stringify({done: true}));
    });
}

function score() {
    d3.select("#score").style("display", "block");
    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        fsm.done();
    });
}

// % operator in javascript is remainder and isn't helpful for wrapping negatives
function mod(a, n) {
    return a - (n * Math.floor(a/n));
}
