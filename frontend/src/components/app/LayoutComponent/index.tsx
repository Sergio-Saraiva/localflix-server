import { FC } from "react";
import { SideBarComponent } from "../SidebarComponent";
import { Outlet } from "react-router-dom";

export const LayoutComponent: FC = () => {
    return (
        <div className="flex h-screen">
            <SideBarComponent />
            <div className="flex-1 overflow-y-auto p-8">
                <Outlet />
            </div>
        </div>  
      )
}