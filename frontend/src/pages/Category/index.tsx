import { MediaLibraryComponent } from "@/components/app/MediaLibraryComponent";
import { FC, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { GetCategory,  ListFolderByCategory } from "../../../wailsjs/go/main/App";
import { models } from "../../../wailsjs/go/models";

export const CategoryPage: FC = () => {
    const [category, setCategory] = useState<models.Category>();
    const [folders, setFolders] = useState<models.Folder[]>([]);
    const { id } = useParams<{ id: string }>();

    const fetchCategory = async () => {
        const result = await GetCategory(parseInt(id!));
        console.log(result)
        setCategory(result);
    }

    const fetchFoldersByCategory = async (id: number) => {
        const result = await ListFolderByCategory(id)
        console.log(result)
        setFolders(result)
    }

    useEffect(() => {
        fetchCategory()
        fetchFoldersByCategory(parseInt(id!))
    }, [id])

    return <MediaLibraryComponent title={category?.Name!} categoryId={category?.ID!} folders={folders} />
}