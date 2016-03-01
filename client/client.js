// hardcoded for testing
var game = 0;
var player = 0;

var probColor = d3.scale.quantize().domain([1, 100]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var costColor = d3.scale.quantize().domain([10, 250]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var taxColor = d3.scale.quantize().domain([100, 550]).range(["#4575b4","#91bfdb","#e0f3f8","#fee090","#fc8d59","#d73027"]);
var oilColor = d3.scale.quantize().domain([1, 9]).range(["#4d4d4d",,"#878787","#bababa","#e0e0e0","#ffffff","#fddbc7","#f4a582","#d6604d","#b2182b"]);

var state = {};

var fsm = StateMachine.create({
    initial: 'lobby',

    events: [
        { name: 'done', from: 'lobby', to: 'survey' },
        { name: 'done', from: 'survey', to: 'report' },
        { name: 'yes', from: 'report', to: 'drill' },
        { name: 'done', from: 'report', to: 'wells' },
        { name: 'done', from: 'drill', to: 'wells' },
        { name: 'done', from: 'wells', to: 'score' },
        { name: 'done', from: 'score', to: 'survey' },
        { name: 'survey', from: 'lobby', to: 'survey' },
        { name: 'report', from: 'lobby', to: 'report' },
        { name: 'drill', from: 'lobby', to: 'drill' },
        { name: 'wells', from: 'lobby', to: 'wells' },
        { name: 'score', from: 'lobby', to: 'score' },
    ],

    callbacks: {

        onenterlobby: lobby,
        onentersurvey: survey,
        onenterreport: report,
        onenterdrill: drill,
        onenterwells: wells,
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
        onleavewells: function() {
            d3.select("#wells").style("display", "none");
            d3.select("#wells-table tbody").html("");
            Mousetrap.reset();
        },
        onleavescore: function() {
            d3.select("#score").style("display", "none");
            Mousetrap.reset();
        },
    }
});

function moveURL() {
    return "/game/" + game + '/player/' + player + '/';
}

function lobby() {
    updateState();

    function updateState() {
        var poll;
        d3.json("/game/" + game + "/")
            .on("load", function(update) {
                if (!update.started) {
                    d3.select("#lobby").style("display", "block");
                }

                d3.select("#lobby-players")
                    .selectAll("tr")
                    .data(update.players)
                    .enter()
                    .append("tr")
                    .selectAll("td")
                    .data(function(d) { return d; })
                    .enter()
                    .append("td")
                    .text(function(d) { return d; });

                if (update.started) {
                    clearTimeout(poll);

                    d3.json(moveURL())
                        .on("load", function(data) {
                            state = data;
                            fsm[state.name]()
                        })
                        .on("error", function(error) {
                            console.log(error);
                        })
                        .get();
                }
            })
            .on("error", function(error) {
                console.log(error);
            })
            .get();
        poll = setTimeout(updateState, 1000);
    }

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json(moveURL())
            .on("load", function(data) {} )
            .on("error", function(error) {
                console.log(error);
            })
            .post(JSON.stringify(-1));
    });
}

function toCurrency(cents, width) {
    s = cents + '';
    s = s.length >= 3 ? s : new Array(3 - s.length + 1).join(0) + s;
    return "$" + s.slice(0, -2) + "." + s.slice(-2)
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
    d3.select("#week").text("Week " + state.week);
    d3.selectAll("#survey-price").text(toCurrency(state.price));
    d3.selectAll("#survey-week").text(state.week);

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

        d3.json(moveURL())
            .on("load", function(data) {
                state = data;
                if (state.name != 'survey') {
                    fsm.done();
                }
            })
            .on("error", function(error) {
                console.log(error);
             })
            .post(JSON.stringify(site));
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
    d3.select("#report-prob").text(state.prob + "%");
    d3.select("#report-cost").text("$\t" + state.cost);
    d3.select("#report-tax").text("$\t" + state.tax);

    Mousetrap.bind('y', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json(moveURL())
            .on("load", function(data) {
                state = data;
                fsm.yes();
            })
            .on("error", function(error) { console.log(error); })
            .post(JSON.stringify(1));
    });
    Mousetrap.bind('n', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json(moveURL())
            .on("load", function(data) {
                state = data;
                fsm.done();
            })
            .on("error", function(error) { console.log(error); })
            .post(JSON.stringify(0));
    });
}

function drill() {
    advance();

    d3.select("#drill").style("display", "block");

    var site = state.site;
    var bit = 0;
    function advance() {
        d3.json(moveURL())
            .on("load", function(data) {
                state = data;
                if (state.name != 'drill') {
                    return fsm.done()
                }
                bit++;

                d3.select("#drill-depth").text(state.depth);
                d3.select("#drill-cost").text(state.cost);
            })
            .on("error", function(error) {
                console.log(error);
             })
            .post(JSON.stringify(1));
    }

    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        advance();
    });

    Mousetrap.bind('q', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json(moveURL())
            .on("load", function(data) {
                state = data;
                fsm.done();
            })
            .on("error", function(error) { console.log(error); })
            .post(JSON.stringify(-1));
    });
}

function wells() {
    d3.select("#wells-player").text(state.player)
    d3.select("#wells-price").text(toCurrency(state.price))
    d3.select("#wells-week").text(state.week)
    d3.select("#wells").style("display", "block");

    function siteData(d) {
        var x = mod(d.site, 80);
        var y = Math.floor(d.site/80);
        var data = [x, y, d.depth, "$", d.cost, "$", d.tax, "$", d.income, "$", d.pnl]
        return data;
    }

    d3.select("#wells-table tbody")
        .selectAll("tr")
        .data(state.wells, function (d) { return d.week } )
        .enter()
        .append("tr")
        .selectAll("td")
        .data(siteData)
        .enter()
        .append("td")
        .text(function(d) { return d; })
        .order();

    Mousetrap.bind('q', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json(moveURL())
            .on("load", function(data) {
                state = data;
                if (state.name != 'wells') {
                    fsm.done();
                }
            })
            .on("error", function(error) { alert(error); })
            .post(JSON.stringify(-1));
    });
}

function score() {
    d3.select("#score").style("display", "block");
    Mousetrap.bind('space', function(e) {
        e.preventDefault ? e.preventDefault() : (e.returnValue = false);
        d3.json(moveURL())
            .on("load", function(data) {
                state = data;
                if (state.name != 'score') {
                    fsm.done();
                }
            })
            .on("error", function(error) { alert(error); })
            .post(JSON.stringify(-1));
    });
}

// % operator in javascript is remainder and isn't helpful for wrapping negatives
function mod(a, n) {
    return a - (n * Math.floor(a/n));
}
