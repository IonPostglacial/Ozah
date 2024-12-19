class HazoTreeMenu extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        this.addEventListener("click", (e) => {
            const selectedItem = e.target.closest("li");
            const items = this.querySelectorAll("li.selected");
            items.forEach(item => item.classList.remove("selected"));
            selectedItem.classList.add("selected");

            const toggleButton = e.target.closest(".btn.toggle");
            if (toggleButton !== null) {
                toggleButton.classList.toggle("primary");
            }
        });
    }
}

customElements.define("hazo-treemenu", HazoTreeMenu);