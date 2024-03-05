(() => {
    function onLoad() {
        const pictureLink = document.getElementById("pictureLink");
        const pictureDialog = document.getElementById("pictureDialog");
        const pictureDialogCloseBtn = document.getElementById("pictureDialogCloseBtn");
        
        function openPictureDialog() {
            pictureDialog.showModal();
        }
    
        function closeDialog(e) {
            e.preventDefault();
            pictureDialog.close();
        }

        pictureLink.addEventListener("click", openPictureDialog);
        pictureDialogCloseBtn.addEventListener("click", closeDialog);
    }
    window.addEventListener("load", onLoad);
})();