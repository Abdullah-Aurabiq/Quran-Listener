window.onload = function () {
    console.log("Window is loaded");
    const logo = document.querySelector("body > div > spline-viewer")?.shadowRoot?.querySelector("#logo");
    console.log("Logo is", typeof logo);
    console.log("Logo is", logo)
    if (logo != null) {
        console.log("Logo is not Null!");
        // logo.style.display = "none";
        logo.remove()
        console.log("Logo has been removed");

    } else {
        setTimeout(() => {
            const logo = document.querySelector("body > div > spline-viewer")?.shadowRoot?.querySelector("#logo");
            logo.style.display = "none";
        }, 1000)
        console.log("Logo not found in 1 second, trying again...");
    }
}