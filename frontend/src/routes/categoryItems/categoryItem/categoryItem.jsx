import style from "./categoryItem.module.css"

export default function CategoryItem({name, amount, image}) {
    return(
        <>
            <article className={style.category}>
                <img src={image} alt={name}></img>
                <h3>{name}</h3>
                <p>{amount + ` articles`}</p>
            </article>
        </>
    )
}