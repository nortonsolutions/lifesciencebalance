/**
 * HomeView
 * @author Norton 2022
 */
 class NavView {

    constructor(router) {
        this.router = router;
    }

    async render(context) {

        context.indexActive = 'active';

        if (!this.navTemplate) {
            this.navTemplate = await getTemplate('app/views/templates/navbar.hbs');
        } 

        if (!this.footerTemplate) {
            this.footerTemplate = await getTemplate('app/views/templates/footer.hbs');
        } 

        this.nav = document.getElementById('nav');
        this.nav.innerHTML = this.navTemplate(context);

        this.foot = document.getElementById('foot');
        this.foot.innerHTML = this.footerTemplate(context);

        this.router.setRouteLinks(this.nav, this.foot);

        const navLinks = document.querySelectorAll('.nav-item:not(.dropdown), .dropdown-item');
        const menuToggle = document.getElementById('navbarNavDropdown');
        const bsCollapse = new bootstrap.Collapse(menuToggle, {toggle:false});

        navLinks.forEach((l) => {
            l.addEventListener('click', () => { 
                if (menuToggle.classList.contains('show')) bsCollapse.toggle();
            })
        })

        document.title = `nortonApp`;
        window.scrollTo(0,0);
    }

}

export { NavView };