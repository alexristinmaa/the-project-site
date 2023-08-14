import {
    useLoaderData
} from "react-router-dom";

import { getPosts } from "../scripts/posts";

import PostListItem from "./postListItem/postListItem";
import PostListSelector from "./postListSelector/postListSelector";
import CategoryItems from "./categoryItems/categoryItems";

import style from "./root.module.css"

export async function loader() {
  const posts = await getPosts();
  return { posts };
}

export default function Root() {
    const {posts} = useLoaderData();

    function menuClick(e) {
        if(e.target.id != style.menuButton) {
            // So the lines don't rotate individually...
            e.target.parentElement.classList.toggle(style.active);
        } else {
            e.target.classList.toggle(style.active);
        }

        document.getElementsByTagName("nav")[0].classList.toggle(style.active);
    }

    return(
        <>
            <div id={style.root}>
                <header id={style.header}>
                    <img src="logo.png" alt="The Red Parrot in my mind" className={style.logo} />
                    <div id={style.navContainer}>
                        <div id={style.menuButton} onClick={menuClick}>
                            <div></div>
                            <div></div>
                            <div></div>
                        </div>
                        <nav>
                            <span>What is this place?</span>
                            <span>Who am I?</span>
                        </nav>
                    </div>
                </header>
                <main id={style.main}>
                    <article id={style.featured}>
                        <div id={style.featuredImage}>
                            <img src="test-thumbnail.jpeg" alt="A Test Thumbail"/>
                        </div>
                        <div id={style.featuredText}>
                            <p className={style.date}>August 9, 2023</p>
                            <h2 className={style.header}>Injecting XXS Payloads into Web Applications using Three.js</h2>
                            <p className={style.description}>Bypassing specially designed barriers using secret tricks of the trade. Follow me on a research journey into MacOS aliases.</p>
                            <div className={style.footer}>
                                <div className={style.divisor}></div>
                                <div className={style.tag}>Hacking</div>
                            </div>
                        </div>
                    </article>
                    <h2 id={style.latestHeader}>Latest Posts</h2>
                    <div id={style.latestList}>
                        <PostListItem post={{Date:"9 August 2023", Title:"Injecting XXS Payloads into Web Applications using Three.js", Description:"Bypassing specially designed barriers using secret tricks of the trade. Follow me on a research journey into MacOS aliases.", Tag:"Hacking"}}/>
                        <PostListItem post={{Date:"9 August 2023", Title:"Injecting XXS Payloads into Web Applications using Three.js", Description:"Bypassing specially designed barriers using secret tricks of the trade. Follow me on a research journey into MacOS aliases.", Tag:"Hacking"}}/>
                        <PostListItem post={{Date:"9 August 2023", Title:"Injecting XXS Payloads into Web Applications using Three.js", Description:"Bypassing specially designed barriers using secret tricks of the trade. Follow me on a research journey into MacOS aliases.", Tag:"Hacking"}}/>
                    </div>
                    <div id={style.selector}>
                        <PostListSelector length={posts.length} current={1}/>   
                    </div>

                    <div id={style.search}>
                        <h2 id={style.looking}>Looking for something <span>specific?</span></h2>

                        <div id={style.searchBar}>
                            <i class="search-icon"></i>
                            <input placeholder="Search the blog..."></input>
                        </div>
                    </div>
                    <h2 id={style.categoriesHeader}><span>Or</span> browse categories.</h2>
                    <div id={style.categories}>
                        <CategoryItems categories={[{
                            name: "Hacking",
                            amount: 23,
                            image: "test-thumbnail.jpeg"
                        },
                        {
                            name: "Cooking",
                            amount: 42,
                            image: "test-thumbnail.jpeg"
                        },
                        {
                            name: "Programming",
                            amount: 12,
                            image: "test-thumbnail.jpeg"
                        },
                        {
                            name: "Life",
                            amount: 9,
                            image: "test-thumbnail.jpeg"
                        },
                        {
                            name: "Gaming",
                            amount: 3,
                            image: "test-thumbnail.jpeg"
                        },
                        {
                            name: "Building",
                            amount: 52,
                            image: "test-thumbnail.jpeg"
                        },
                        {
                            name: "Photography",
                            amount: 61,
                            image: "test-thumbnail.jpeg"
                        },
                        {
                            name: "Fun",
                            amount: 1,
                            image: "test-thumbnail.jpeg"
                        }
                        ]}
                        />
                    </div>
                </main>
                <footer>
                    <span>Alexander Ristinmaa, 2023</span>
                    <span>Made with React, Go and <svg viewBox="0 0 1792 1792" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg" class="inline-icon"><path d="M896 1664q-26 0-44-18l-624-602q-10-8-27.5-26T145 952.5 77 855 23.5 734 0 596q0-220 127-344t351-124q62 0 126.5 21.5t120 58T820 276t76 68q36-36 76-68t95.5-68.5 120-58T1314 128q224 0 351 124t127 344q0 221-229 450l-623 600q-18 18-44 18z" fill="#e25555"></path></svg></span>
                    <span>Creation documented <a>here</a></span>
            </footer>
            </div>
        </>
    )
}