(function (){
    const toggles = document.querySelectorAll(".textToggle");
    console.log("hello.")
    toggles.forEach(t => {
        t.addEventListener("click", function() {
            t.classList.toggle("isToggled");
        })
    })

})()
