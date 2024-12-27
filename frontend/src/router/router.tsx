import { LayoutComponent } from "@/components/app/LayoutComponent";
import { HomePage } from "@/pages/Home";
import { CategoryPage } from "@/pages/Category";
import { SettingsPage } from "@/pages/Settings";
import { createBrowserRouter } from "react-router-dom";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <LayoutComponent />,
        children: [
            {
                path: '',
                element: <HomePage />
            },
            {
                path: 'settings',
                element: <SettingsPage />
            },
            {
                path: 'category/:id',
                element: <CategoryPage />
            }
        ]
        
    }
])