function onLoad() {
    document.body.addEventListener('click', function (e) {
        const menu = e.target.closest("menu");
        if (menu !== null) {
            const selectedItem = e.target.closest("li");
            const items = menu.querySelectorAll("li.selected");
            items.forEach(item => item.classList.remove("selected"));
            selectedItem.classList.add("selected");
        }
        const toggleButton = e.target.closest(".btn.toggle");
        if (toggleButton !== null) {
            toggleButton.classList.toggle("primary");
        }
    });
}
window.addEventListener("load", onLoad);