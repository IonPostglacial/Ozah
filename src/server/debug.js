(function() {
    let lastStartedOn;
    setInterval(async function () {
        const startedOn = await fetch("/started-on");
        const startedOnTxt = await startedOn.text();
        if (typeof lastStartedOn !== "undefined" && typeof startedOnTxt !== "undefined" && startedOnTxt !== lastStartedOn) {
            location.reload();
        } else {
            lastStartedOn = startedOnTxt;
        }
    }, 3000);
})();