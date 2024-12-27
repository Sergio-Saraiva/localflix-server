import { Film, Home, Settings } from "lucide-react";
import { FC, useEffect, useState } from "react";
import { NavLink } from "react-router-dom";
import { ListCategories } from "../../../../wailsjs/go/main/App";
import { models } from "../../../../wailsjs/go/models";

export const SideBarComponent: FC = () => {
    const [categories, setCategories] = useState<models.Category[]>([]);
    const fetchCategories = async () => {
        const result = await ListCategories();
        console.log(result)
        setCategories(result);
    }

    useEffect(() => {
        fetchCategories()
    }, [])

    return (
        <div className="w-64 bg-white h-full shadow-md">
          <div className="p-6">
            <h1 className="text-2xl font-bold text-gray-800">LocalFlix</h1>
          </div>
          <nav className="mt-6">
            <NavLink to="/" className="block py-2 px-6 hover:bg-gray-100">
              <Home className="inline-block mr-2" size={20} />
              Home
            </NavLink>
            {
                categories.map((category, index) => (
                    <NavLink key={index} to={`/category/${category.ID}`} className="block py-2 px-6 hover:bg-gray-100">
                        <Film className="inline-block mr-2" size={20} />
                        {category.Name}
                    </NavLink>
                ))
            }
            <NavLink to="/settings" className="block py-2 px-6 hover:bg-gray-100">
              <Settings className="inline-block mr-2" size={20} />
              Settings
            </NavLink>
          </nav>
        </div>
    )
}