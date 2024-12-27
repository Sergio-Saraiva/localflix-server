import { Button } from "@/components/ui/button";
import { Folder, Trash2 } from "lucide-react";
import { FC } from "react";

interface FolderItemProps {
    path: string
}

export const FolderItemComponent: FC<FolderItemProps> = ({ path }) => {
    return (
        <li className="flex items-center justify-between bg-gray-100 p-3 rounded">
          <div className="flex items-center">
            <Folder className="mr-2" size={20} />
            <span>{path}</span>
          </div>
          <Button variant="destructive" size="icon">
            <Trash2 size={16} />
          </Button>
        </li>
      )
}