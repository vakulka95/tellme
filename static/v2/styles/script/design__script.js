

// Show and hide modal
const filter_btn = document.querySelector('.filter .filter-btn');
const close_btn = document.querySelector('.close');
const modal = document.querySelector('.modal');
const body = document.querySelector('body');
const user = document.querySelector('.user');
const userImg = document.querySelector('.user img');


filter_btn.onclick = showModal;

function showModal(){
    modal.classList.add('show');
    body.style.overflowY = 'hidden';
}

close_btn.onclick = closeModal;

function closeModal(){
    modal.classList.remove('show');
    body.style.overflowY = 'scroll';
}

// Active nav link
const menuLink= document.querySelectorAll('.menu--link');

menuLink.forEach(el => {
    el.onclick = () => {
        if(!el.classList.contains('menu--link--active')){
            menuLink.forEach(link => link.classList.remove('menu--link--active'));
            el.classList.add('menu--link--active');
        }
    }
})

// Hide elements

const filter = document.querySelector('.filter');
const filterDesk = document.querySelector('.filter-desktop');
const appClients = document.querySelector('.app__clients');
const desktopNav = document.querySelector('.desktop-nav');
const psychologists = document.querySelector('.psychologists');
const header = document.querySelector('.header');
const main = document.querySelector('.main');
const rate = document.querySelector('.rate');
const feedbacks = document.querySelector('.feedbacks');



function hideFilterBlock(){
    if(document.body.clientWidth >= 1024){
        filter.style.display = 'none';
        filterDesk.style.display = 'block';
        appClients.style.display = 'none';
        header.style.position = 'relative';
        main.style.display = 'block';
        userImg.style.display = 'none';
        user.style.display = 'none';
        desktopNav.style.display = 'block';
        psychologists.style.display = 'none';
        feedbacks.style.display = 'none';

    }else{
        filter.style.display = 'block';
        filterDesk.style.display = 'none';
        appClients.style.display = 'block';
        header.style.position = 'sticky';
        main.style.display = 'none';
        userImg.style.display = 'block';
        desktopNav.style.display = 'none';
        psychologists.style.display = 'block';
        feedbacks.style.display = 'block';


    }
}
hideFilterBlock();

window.addEventListener('resize', hideFilterBlock);

// Background table

const tablesTr = document.querySelectorAll('table tr');

function spillTr(){
    for(let i = 1; i < tablesTr.length; i++){
        if(i % 2 === 0 ){
            tablesTr[i].style.backgroundColor = '#E3E8F5';
        }
    }
}

spillTr();

// SEARCH NAME 
const filterSearch = document.querySelector('.filter-desktop-search');
const nameInput = document.querySelectorAll('.name');
const phoneInput = document.querySelectorAll('.phone');

filterSearch.oninput = function(){
    let val = this.value.trim();
    console.log(val)

    if(val != ''){
        nameInput.forEach(el => {
            if(el.innerText.search(val) == -1 ){
                el.parentElement.parentElement.classList.add('hide')
            }
            else{
                el.parentElement.parentElement.classList.remove('hide');
            }
        });
        phoneInput.forEach(el => {
            if(el.innerText.search(val) == -1 ){
                el.parentElement.parentElement.classList.add('hide')
            }
            else{
                el.parentElement.parentElement.classList.remove('hide');
            }
        });
    }else{
        nameInput.forEach(el => {
            el.parentElement.parentElement.classList.remove('hide')
        });
        phoneInput.forEach(el => {
            el.parentElement.parentElement.classList.remove('hide')
        });
    }
}




// change Arrow
const titleArrow = document.querySelector('.title__arrow');
const form = document.querySelector('.form');
const button = document.querySelector('.button');
const arrow = document.querySelector('.arrow');
const userImg = document.querySelector('.user img');
const user = document.querySelector('.user');
const desktopNav = document.querySelector('.desktop-nav')
const menuModal = document.querySelector('.menu-modal');
const gender = document.querySelector('.app-gender-form');
const form = document.querySelector('.form');
const arrow = document.querySelector('.arrow');

function changeArrow(){
    if(document.body.clientWidth >= 1024){
        form.style.width = '75%';
        userImg.style.display = 'none';
        user.style.display = 'none';
        desktopNav.style.display = 'flex';
        gender.style.display = 'none';

    }else{
        form.style.width = '100%';
        userImg.style.display = 'block';
        user.style.display = 'block';
        desktopNav.style.display = 'none';
        gender.style.display = 'block';



    }
}
changeArrow();

window.addEventListener('resize', changeArrow); 


user.onclick = () => {

    menuModal.classList.toggle('show');
    menuModal.onclick = () => { 
    menuModal.classList.remove('show');
    }
}

// Active nav link
const user = document.querySelector('.user');
const userImg = document.querySelector('.user img');
const menuModal = document.querySelector('.menu-modal');
const drag = document.querySelector('.drag');
const dragTitle = document.querySelector('.drag-title');


// show modal menu
user.onclick = () => {

    menuModal.classList.toggle('show');
    menuModal.onclick = () => { 
    menuModal.classList.remove('show');
    }
}

const desktopNav = document.querySelector('.desktop-nav');


function hideBlock(){
    if(document.body.clientWidth >= 1024){
        drag.classList.add('d-flex');
        dragTitle.style.display = 'block';
        user.style.display = 'none';
        userImg.style.display = 'none';
        desktopNav.style.display = 'block';
    }else{
        drag.classList.remove('d-flex')
        drag.style.display = 'none';
        dragTitle.style.display = 'none';
        user.style.display = 'block';
        userImg.style.display = 'block';
        desktopNav.style.display = 'none';

    }
}
hideBlock();

window.addEventListener('resize', hideBlock);


// Edit inputs





// Search application in search input


const filterSearch = document.querySelector('.filter-desktop-search');
const nameInput = document.querySelectorAll('.name');
const phoneInput = document.querySelectorAll('.phone');

    filterSearch.oninput = function(){
        let val = this.value.trim();
        console.log(val)

        if(val != ''){
            nameInput.forEach(el => {
                if(el.innerText.search(val) == -1 ){
                    console.log( el.parentElement.parentElement);
                    el.parentElement.parentElement.parentElement.classList.add('hide')
                }
                else{
                    el.parentElement.parentElement.parentElement.classList.remove('hide');
                }
            });
        
        }else{
            nameInput.forEach(el => {
                el.parentElement.parentElement.parentElement.classList.remove('hide')
            });
            
        }
    }

const filterMobSearch = document.querySelector('.filter-mob-input');
const nameMob = document.querySelectorAll('.name-mob');

filterMobSearch.oninput = function(){
    let val = this.value.trim();
    if(val != ''){
        nameMob.forEach(el => {
            if(el.innerText.search(val) == -1 ){
                el.parentElement.parentElement.parentElement.parentElement.parentElement.classList.add('hide')
            }
            else{
                el.parentElement.parentElement.parentElement.parentElement.parentElement.classList.remove('hide');
            }
        });
        
    }else{
        nameMob.forEach(el => {
            el.parentElement.parentElement.parentElement.parentElement.parentElement.classList.remove('hide')
        });
        
    }
}


