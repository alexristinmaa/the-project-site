// Import the functions you need from the SDKs you need

import { initializeApp } from "firebase/app";

import { getAnalytics } from "firebase/analytics";

import { getFirestore, collection, getDocs } from "firebase/firestore"; 

// TODO: Add SDKs for Firebase products that you want to use

// https://firebase.google.com/docs/web/setup#available-libraries


// Your web app's Firebase configuration

// For Firebase JS SDK v7.20.0 and later, measurementId is optional

const firebaseConfig = {
    apiKey: "AIzaSyBJjPcpRGWgh2SzsMN94Z1TN5zH68vV1Nc",
    authDomain: "the-project-site.firebaseapp.com",
    projectId: "the-project-site",
    storageBucket: "the-project-site.appspot.com",
    messagingSenderId: "929843183431",
    appId: "1:929843183431:web:13e4aecd848fa8fb719662",
    measurementId: "G-4VFYK21YNT"
};


// Initialize Firebase

const app = initializeApp(firebaseConfig);

const analytics = getAnalytics(app);
export const db = getFirestore(app);


export async function getPost({id}) {

    let data = {
        id: id
    }

    let response = await fetch(`/api/post`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });

    let post = await response.json();

    return post;
}

export async function getFeaturedPost() {
    console.log(db);
    return {Date:"9 August 2023", Title:"Injecting XXS Payloads into Web Applications using Three.js", Description:"Bypassing specially designed barriers using secret tricks of the trade. Follow me on a research journey into MacOS aliases.", Tag:"hacking", Id:"bypassing-specially-designed-barriers"};
}

export async function getPosts({search, tags, page}) {
    const postsCol = collection(db, 'posts');
    const postSnapshot = await getDocs(postsCol);
    const postList = postSnapshot.docs.map(doc => doc.data());

    console.log(postList);
    return {posts: postList, total: postList.length};


    /*
    page ||= 1;
    tags ||= [];
    search ||= "";

    let data = {
        page: page,
        tags: tags,
        search: search
    }

    let response = await fetch("/api/posts/", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });

    let posts = await response.json();

    return posts;*/
}

export async function getTags() {
    let response = await fetch("/api/getTags");

    let jsonD = await response.json();

    return jsonD;
} 