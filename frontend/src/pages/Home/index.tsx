import { Input } from "@/components/ui/input";
import { Search } from "lucide-react";
import { FC, useEffect, useState } from "react";
import { MediaCard } from "./components";
import { ListCategories } from "../../../wailsjs/go/main/App";
import { models } from "wailsjs/go/models";

export const HomePage : FC = () => {

  const [categories, setCategories] = useState<models.Category[]>([]);

  const fetchCategories = async () => {
    const result = await ListCategories();
    setCategories(result);
  }

    useEffect(() => {
      fetchCategories()
    }, [])

    return (
        <div>
          <h1 className="text-3xl font-bold mb-6">Welcome to LocalFlix</h1>
          <div className="mb-8">
            <div className="relative">
              <Input 
                type="text" 
                placeholder="Search for movies, TV shows, or anime..." 
                className="pl-10 pr-4 py-2 w-full"
              />
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
            </div>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {categories.map((category, index) => (
              <MediaCard key={index} title={category.Name} count={10} link={`/category/${category.ID}`} />
            ))}
          </div>
        </div>
      )
}