import style from "./categoryItem.module.css";

import {
    Link
} from "react-router-dom";

export default function CategoryItem({name, amount, image}) {
    return(
        <>
            <Link to={"/posts/page/1?tags=" + name}>
                <article className={style.category}>
                    <img src={image} alt={name}></img>
                    <h3>{name}</h3>
                    <p>{amount + ` articles`}</p>
                </article>
            </Link>
            
        </>
    )
}