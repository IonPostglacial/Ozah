class HazoPopover extends HTMLElement {
    constructor() {
        super();
    }
}

function onLoad() {
    document.body.addEventListener('click', function (e) {
        console.log("click");
        const clickedPopover = e.target.closest("hazo-popover");
        let selectedCheckbox = null
        if (clickedPopover !== null) {
            const checkbox = clickedPopover.querySelector("input[type=checkbox]");
            if (checkbox.checked) {
                selectedCheckbox = checkbox;
            }
        }
        document.querySelectorAll("hazo-popover>label>input[type=checkbox]").forEach(cb => {
            if (cb !== selectedCheckbox) {
                cb.checked = false;
            }
        });
    });
}

customElements.define("hazo-popover", HazoPopover);
window.addEventListener("load", onLoad);