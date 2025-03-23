class HazoSplitPanelGutter extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        const leftPane = this.parentElement.parentElement.firstElementChild;

        this.addEventListener("mousedown", (e) => {
            e.preventDefault();
            window.addEventListener('mousemove', mousemove);
            window.addEventListener('mouseup', mouseup);
            
            let startX = e.x;
            const startWidth = leftPane.getBoundingClientRect().width;
            
            console.log("start x", startX);
            console.log("start width", startWidth);

            function mousemove(e) {
                const delta = e.x - startX;
                console.log("delta", delta);
                const width = startWidth + delta + "px";
                console.log("new width", width);
                leftPane.style.width = width;
            }
            
            function mouseup() {
              window.removeEventListener('mousemove', mousemove);
              window.removeEventListener('mouseup', mouseup);
              
            }
        });
    }
}

customElements.define("hazo-splitpane-gutter", HazoSplitPanelGutter);