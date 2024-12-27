import { Button } from "@/components/ui/button";
import { FC } from "react";

interface MediaCardProps {
    title: string;
    count: number;
    link: string;
}

export const MediaCard: FC<MediaCardProps> = ({ title, count, link }) => {
    return (
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-2">{title}</h2>
          <p className="text-gray-600 mb-4">{count} items</p>
          <Button asChild>
            <a href={link}>View All</a>
          </Button>
        </div>
      )
}