import style from "./featuredPost.module.css"

import {
    Link
} from "react-router-dom";

export function FeaturedPost({post}) {
    return (
        <>
            <article id={style.post}>
                <Link to={"/posts/" + post.Id}>
                    <div id={style.featuredImage}>
                        <img src="/test-thumbnail.jpeg" alt="A Test Thumbail"/>
                    </div>
                </Link>
                <Link to={"/posts/" + post.Id}>
                    <div id={style.featuredText}>
                        <p className={style.date}>{post.Date}</p>
                        <h2 className={style.header}>{post.Title}</h2>
                        <p className={style.description}>{post.Description}</p>
                        <div className={style.footer}>
                            <div className={style.divisor}></div>
                            <Link to={"/posts/page/1?tags=" + post.Tag}><div className={style.tag}>Hacking</div></Link>
                        </div>
                    </div>
                </Link>
            </article>
        </>
    )
}