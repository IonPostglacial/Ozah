function onLoad() {
    document.body.addEventListener('click', function (e) {
        const clickedPopover = e.target.closest(".popover");
        let selectedPopover = null;
        if (clickedPopover !== null) {
            const checkbox = clickedPopover.querySelector("input[type=checkbox]");
            if (checkbox.checked) {
                selectedPopover = clickedPopover;
            }
        }
        document.querySelectorAll(".popover").forEach(popover => {
            if (popover !== selectedPopover) {
                popover.querySelector("input[type=checkbox]").checked = false;
            }
        });
    });
}
window.addEventListener("load", onLoad);