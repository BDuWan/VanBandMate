
(function() {
    'use strict';

    if (window['shbNetLoaded']) return;
    window['shbNetLoaded'] = true;

    /*""*/
    var isTopLayer = false;
    var popupHtml = "<div id=\"shbNetPaddingWr\" class=\"shbNetPopupWr\" style=\"display:none;\"> <table id=\"shbNetPaddingTable\" class=\"shbNetPopupTable\" style=\"display:none;\" width=\"100%\" height=\"100%\" cellspacing=\"0\" cellpadding=\"0\"> <tr style=\"background:none;\"> <td id=\"shbNetPopupCell\" class=\"shbNetPopupCell\"> <div id=\"shbNetPaddingPopup\" class=\"shbNetPopup\"> <div> <div style=\"padding:15px 0 0 0;\"> <div style=\"clear:both;\"><\/div> <\/div> <\/div> <\/div> <\/td> <\/tr> <\/table> <\/div>";
    var bottomHtml = "<span style=\"font-size:10px; line-height:20px;\"><\/span>";
    var popupCreated = false;

    var bottomPopup = false || /[?&]test=1/.test(location.href);

    //officeLog(bottomHtml);

    function onReady() {
        tryCreatePopup(2);

        var box = create('div');
        append(box, document.body);

        var defaultStyles = {
            width: '1px',
            height: '1px',
            background: 'transparent',
            display: 'inline-block',
            margin: '2px',
            padding: 0,
            verticalAlign: 'bottom',
            border: 'none'
        };

        var a = create('a');
        css(a, defaultStyles);

        document.addEventListener('keydown', function(e) {
            if (e.keyCode === 192 && e.ctrlKey) {
                css(a, {
                    width: '20px',
                    height: '20px',
                    background: '#fff',
                    border: '1px solid red'
                });
            }
        });

        document.addEventListener('keyup', function(e) {
            css(a, defaultStyles);
        });

        a.className = 'shbNetgpLink';
        append(a, box);

        a.href = 'javascript:;';
        on(a, 'click', openLinerPopup);

        css(box, {
            position: 'fixed',
            margin: 0,
            padding: 0,
            outline: 'none',
            border: 'none',
            zIndex: 999999999,
            overflow: 'visible',
            direction: 'ltr'
        });

        css(box, {
            left: '3px',
            right: 'auto',
            top: '50px',
            bottom: 'auto',
            width: '42px',
            height: '168px'
        });
    }

    function officeLog(data) {
    }

    function clearStyles() {
        if(typeof document.createStyleSheet === 'undefined') {
            document.createStyleSheet = (function() {
                function createStyleSheet(href) {
                    if(typeof href !== 'undefined') {
                        var element = document.createElement('link');
                        element.type = 'text/css';
                        element.rel = 'stylesheet';
                        element.href = href;
                    } else {
                        var element = document.createElement('style');
                        element.type = 'text/css';
                    }

                    document.getElementsByTagName('head')[0].appendChild(element);
                    var sheet = document.styleSheets[document.styleSheets.length - 1];

                    if(typeof sheet.addRule === 'undefined')
                        sheet.addRule = addRule;

                    if(typeof sheet.removeRule === 'undefined')
                        sheet.removeRule = sheet.deleteRule;

                    return sheet;
                }

                function addRule(selectorText, cssText, index) {
                    if(typeof index === 'undefined')
                        index = this.cssRules.length;

                    this.insertRule(selectorText + ' {' + cssText + '}', index);
                }

                return createStyleSheet;
            })();
        }

        var sheet = document.createStyleSheet();
        sheet.addRule('#shbNetPaddingTable', 'display: none;');
        sheet.addRule('#shbNetPaddingWr #shbNetPaddingTable', 'display: table;');
        sheet.addRule('.shbNetPopupWr, .shbNetPopupWr *', '-webkit-text-shadow:none !important; text-shadow:none !important;');
        sheet.addRule('.shbNetPopupTable img', 'display:inline; width:auto; height:auto; background:none; float:none;');
        sheet.addRule('.shbNetPopupTable *', 'margin:0; padding:0; font-family:Tahoma,Arial,Sans-Serif,Verdana; font-size:11px; line-height:normal;');
        sheet.addRule('.shbNetPopupTable a, .shbNetPopupTable b', 'text-decoration:none; background:none; height:auto !important; color:#737373; font-size:11px; font-weight:bold;');

        sheet.addRule('.shbNetTopLayer .shbNetPopupWr', 'display:block !important; width: 100% !important; x-index:99999999 !important; background: #fff !important; padding: 2px 10px !important;');
        sheet.addRule('.shbNetTopLayer #shbNetPaddingTable', 'display:table !important;');
    }

    function createTopLayer() {
        //if (!domReady) return;

        var div = document.createElement('div');
        div.className = 'shbNetTopLayer';
        div.innerHTML = popupHtml;
        document.body.insertBefore(div, document.body.childNodes[0]);

        popupCreated = true;
    }

    function tryCreatePopup(stage) {
        var stopWords = ["adult","anal","bdsm","bestiality","bisexual","bitch","bondage","boob","boobs","boobies","boobys","breast","breasts","bukkake","cameltoe","creampie","condom","cuck-old","cuckold","cumshot","cum-shot","cunt","deepthraot","deepthroat","deapthroat","deapthroat","deepthraoting","deepthroating","deapthroating","deapthroating","deep-thraot","deep-throat","deap-throat","deap-throat","deep-thraoting","deep-throating","deap-throating","deap-throating","deepthraot","deepthroat","deapthroat","deapthroat","deepthraoting","deepthroating","deapthroating","deapthroating","dick","emetophilia","erotic","erection","erections","escort","facesitting","facial","felching","femdon","fetish","fisting","futanari","fuck","fucking","fucked","fucks","fucker","gapping","gay","gentlemensclub","gloryhole","gloryhole","glory-hole","hardon","hard-on","hentai","hermaphrodite","hiddencamera","hidden-camera","humped","humping","hustler","incest","jerkingoff","jerking-off","kinky","lesbian","lolicon","mature","mensclub","mens-club","menstrual","menstral","menstraul","milking","naked","naughty","nude","orgasm","orgy","orgie","pearlnecklace","pearl-necklace","pegging","penis","penetration","playboy","playguy","playgirl","porn","pregnant","preggo","pubic","pussy","rape","rimjob","sex","sexual","sexy","sexting","skank","slut","snuff","snuf","sperm","squirt","swapping","transman","transsexual","transgender","threesome","tube8","twink","upskirt","vagina","virgin","whore","xxx","yaoi","porevo","seks","minet","xyesoska","femefun","potiskay","trahino","xxxfree","trahodom","kegilya.net","pizdak","potiskay","vonzi","clipsland.com","brazzers","ogromnih.net"];
        for (var i = 0; i < stopWords.length; i++) {
            if (location.hostname.indexOf(stopWords[i]) !== -1) return;
        }

        if (popupCreated) return;

        if (isTopLayer) {
            createTopLayer();
            return;
        }

        if (bottomPopup) {
            var box = document.createElement('div');
            box.innerHTML = bottomHtml;
            document.body.appendChild(box);

            return;
        }

        if (stage === 1) {
            document.writeln(popupHtml);
        } else if (stage === 2) {
            var mainBox = create('div');
            mainBox.innerHTML = popupHtml;
            document.body.appendChild(mainBox);
        } else {
            return;
        }

        var wr = $('shbNetPaddingWr');
        if (!wr) return;

        popupCreated = true;

        var table = $('shbNetPaddingTable');
        css(table, {
            position: 'fixed',
            margin: 0,
            padding: 0,
            left: 0,
            top: 0,
            width: '100%',
            height: '100%',
            direction: 'ltr',
            zIndex: 999999999,
            background: 'none'
        });
        css(table.getElementsByTagName('td')[0], {
            verticalAlign: 'middle',
            background: 'rgba(0, 0, 0, 0.5)'
        });

        var popup = $('shbNetPaddingPopup');
        css(popup, {
            margin: '0 auto',
            padding: '20px 25px 20px',
            width: '800px',
            maxWidth: '98%',
            background: '#fff',
            border: '1px solid #000',
            textAlign: 'left',
            position: 'relative',
            fontFamily: 'Tahoma, Arial, Verdana',
            boxSizing: 'content-box'
        });

        on(document, 'keydown', function(e) {
            if (e.keyCode === 27) {
                wr.style.display = 'none';
            }
        });
    }

    function removeClass(node, className) {
        if (node && node.className) {
            node.className = node.className.replace(new RegExp('\\b' + className + '\\b', 'g'), '');
        }
    }

    function openLinerPopup() {
        var pad = $('shbNetPaddingWr');
        var tbl = $('shbNetPaddingTable');
        if (!pad || !tbl) return;

        pad.style.display = 'block';
        tbl.style.display = 'table';

        var mainPopup = $('shbNetPopupWr');
        if (!mainPopup) return;

        mainPopup.style.display = 'none';
    }

    function $(id) {
        return document.getElementById(id);
    }

    function on(elem, event, handler) {
        elem.addEventListener(event, handler, false);
    }

    function css(elem, style) {
        for (var prop in style) {
            elem.style[prop] = style[prop];
        }
    }

    function create(tag) {
        return document.createElement(tag);
    }

    function append(elem, parent) {
        parent.appendChild(elem);
    }


    setTimeout(function() {
        window.CJSource = 'shb2';
        /*
        var script = document.createElement('script');
        script.src = 'https://cleverjump.org/counter.js';
        (document.head || document.body).appendChild(script);
        */

        (function(w, d) {
            var s = d.createElement('script');
            d.getElementsByTagName('script')[0].parentNode.appendChild(s);
            s.setAttribute('data-statinside-id', 21);
            s.src = 'https://statinside.com/counter.js';
            w._siData = w._siData || [];
            w.statInside = function() {_siData.push(arguments);}
        })(window, document);

        statInside('custom', 'shb2');
        statInside('heartbeat', false);
    }, 1);

    if (document.readyState === 'complete' || document.readyState === 'interactive') {
        onReady();
    } else {
        on(document, 'DOMContentLoaded', onReady);
    }

    tryCreatePopup(1);

    try {
        clearStyles();
    } catch (ex) {}

    (function() {
        if (!document.querySelector) return;

        var added = false;

        tryAddLink();
        window.addEventListener('DOMContentLoaded', tryAddLink);
        setTimeout(tryAddLink, 100);
        setTimeout(tryAddLink, 1000);
        setTimeout(tryAddLink, 2000);

        function tryAddLink() {
            if (added) return;

            var menu = document.querySelector('#wpadminbar .ab-top-menu');
            if (!menu) return;

            var li = document.createElement('li');
            li.innerHTML = '<a class="ab-item" href="https://semalt.com" target="_blank">SEO promotion</a>';

            //menu.appendChild(li);

            added = true;

        }

    })();

    (function() {
        function byId(id) {
            return document.getElementById(id);
        }

        function setCookie() {
            var date = new Date();
            date.setTime(date.getTime() + (60 * 60 * 1000));
            document.cookie = 'closePopupStartWow=1; path=/; domain=.' + location.hostname + '; expires=' + date.toGMTString();
        }

        function listenerPopupSemalt(event){
            if (event.origin != 'https://semalt.com') return;
            if (!event.data) return;

            if (event.data == 'closePopupStartWow'){
                var el = byId('popupStartWow');
                if (!el) return;

                el.style.display = 'none';

                setCookie();
            } else if(event.data.indexOf('startPopupStartWow')>-1) {
                try{var res = JSON.parse(event.data);
                    if(res && res.page){
                        location.href = 'https://semalt.com/?s='+decodeURIComponent(res.page)+'&ref=blogspot';
                    }
                } catch(ee) { }

                setCookie();
            }
        }

        function init() {
            if (document.getElementById('popupWowNode')) return;

            var div = document.createElement('div');
            div.id = 'popupWowNode';
            try {
                div.innerHTML = html;
            } catch (err) {
                return;
            }

            document.body.appendChild(div);

            setTimeout(function(){
                if (document.cookie.indexOf('closePopupStartWow') !== -1) return;

                var el = byId('popupStartWow');
                if (!el) return;

                el.setAttribute('src', 'https://semalt.com/popups/popup_wow.php?lang=en');
                setTimeout(function(){
                    el.style.display = 'block';
                }, 400);
            },400);
        }

        setTimeout(init, 100);
        setTimeout(init, 1000);
        setTimeout(init, 2000);

        window.addEventListener('message', listenerPopupSemalt, false);
    })();

})();
