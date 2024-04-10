function onLoad() {
    document.body.addEventListener('click', function (e) {
        const clickedPopover = e.target.closest(".popover");
        let selectedCheckbox = null
        if (clickedPopover !== null) {
            const checkbox = clickedPopover.querySelector("input[type=checkbox]");
            if (checkbox.checked) {
                selectedCheckbox = checkbox;
            }
        }
        document.querySelectorAll(".popover > input[type=checkbox]").forEach(cb => {
            if (cb !== selectedCheckbox) {
                cb.checked = false;
            }
        });
    });
}
window.addEventListener("load", onLoad);