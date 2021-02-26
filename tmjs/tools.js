// @namespace    http://firsh.me/
// @version      1.3.6
// @description  在线看小说 J 下一行  K上一行 C下一章
// @author       一个俗人
// --------------------------------------------------------------------
//
// ==UserScript==
// @name          t_txt
// @namespace     http://diveintogreasemonkey.org/download/
// @description   script to t_txt remote console on every page
// @match         *://google.com/*
// @match         https://www.google.com/search*
// @grant          none
// @run-at         document-start
// @require        https://code.jquery.com/jquery-2.0.3.min.js
// @grant          unsafeWindow
// ==/UserScript==
(function() {
    'use strict';
    console.log("init");
    function request_nl() {
        var settings = {
            "async": true,
            "crossDomain": true,
            "url": "http://127.0.0.1:18311/nl",
            "method": "GET",
            "headers": {
                "cache-control": "no-cache",
            }
        }

        $.ajax(settings).done(function (response) {
            if (response !== ""){
                document.querySelector("#rso > div > div:nth-child(3) > div > div.IsZvec > div > span").innerHTML = response
            }else{
                alert("----")
            }
        });
    }
    function request_pl() {
        var settings = {
            "async": true,
            "crossDomain": true,
            "url": "http://127.0.0.1:18311/pl",
            "method": "GET",
            "headers": {
                "cache-control": "no-cache",
            }
        }

        $.ajax(settings).done(function (response) {
            if (response !== ""){
                document.querySelector("#rso > div > div:nth-child(3) > div > div.IsZvec > div > span").innerHTML = response
            }else{
                alert("----")
            }
        });
    }
    function request_np() {
        var settings = {
            "async": true,
            "crossDomain": true,
            "url": "http://127.0.0.1:18311/np",
            "method": "GET",
            "headers": {
                "cache-control": "no-cache",
            }
        }

        $.ajax(settings).done(function (response) {
            document.querySelector("#rso > div > div:nth-child(3) > div > div.IsZvec > div > span").innerHTML = response
        });
    }




    function myEventHandler(e) {
        var keyCode = e.keyCode;
        console.log(e, keyCode, e.which)
        if(e.key==="j"){
            request_nl()
            return
        }
        if(e.key==="k"){
            request_pl()
            return
        }
        if(e.key==="c"){
            request_np()
            return
        }
    };

    window.addEventListener("keypress", myEventHandler, false);
})();
