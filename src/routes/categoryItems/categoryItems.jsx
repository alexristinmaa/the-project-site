import CategoryItem from "./categoryItem/categoryItem"

export default function CategoryItems({categories}) {
    let categoryElements = [];

    categories.sort((a, b) => b.amount - a.amount)

    for(let category of categories) {
        categoryElements.push(
            <CategoryItem name={category.name} amount={category.amount} image={category.image} />
        )
    }

    return(
        <>
            {categoryElements}
        </>
    )
}