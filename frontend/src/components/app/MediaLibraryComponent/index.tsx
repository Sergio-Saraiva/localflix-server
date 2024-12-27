import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { FC } from "react";
import { models } from "wailsjs/go/models";
import { CreateFolderSource } from "../../../../wailsjs/go/main/App";
import { FolderItemComponent } from "../FolderItemComponent";
interface MediaLibraryProps {
    title: string
    folders: models.Folder[]
    categoryId: number
}

export const MediaLibraryComponent: FC<MediaLibraryProps> = ({ title, folders, categoryId }) => {
    const addFolder = () => {
        console.log('Add Folder')
        CreateFolderSource(categoryId)
    }

    return (
        <div>
          <h1 className="text-3xl font-bold mb-6">{title} Library</h1>
          <div className="bg-white p-6 rounded-lg shadow-md">
            <h2 className="text-xl font-semibold mb-4">Add Folder</h2>
            <div className="flex space-x-2 mb-6">
              <Button onClick={addFolder}>
                <Plus className="mr-2" size={16} />
                Add
              </Button>
            </div>
            <h2 className="text-xl font-semibold mb-4">Current Folders</h2>
            <ul className="space-y-2">
                {folders.map((folder, index) => (
                    <FolderItemComponent key={index} path={folder.path} />
                ))}
            </ul>
          </div>
        </div>
      )
}