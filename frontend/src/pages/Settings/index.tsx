import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { FC, useEffect, useState } from "react";
import { CreateCategory, DeleteCategory, ListCategories, StartServer, StopServer } from "../../../wailsjs/go/main/App";
import { models } from "wailsjs/go/models";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Play, Plus, Square } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

export const SettingsPage : FC = () => {
    const [categories, setCategories] = useState<models.Category[]>([])
    const [newCategory, setNewCategory] = useState('')
    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [isServerRunning, setIsServerRunning] = useState(false)
    const { toast } = useToast()
    
    const fetchCategories = async () => {
        const result = await ListCategories();
        console.log(result)
        setCategories(result);
    }

    useEffect(() => {
        fetchCategories()
    }, [])

    const handleAddCategory = async () => {
        if (newCategory) {
        const category = await CreateCategory(newCategory)
        setCategories([...categories, category])
        setNewCategory('')
        setIsDialogOpen(false)
        }
    }

    const handleRemoveCategory = async (id: number) => {
        await DeleteCategory(id)
        setCategories(categories.filter(category => category.ID !== id))
    }

    const toggleServer = async () => {
        if (isServerRunning) {
        await StopServer()
          toast({ title: "Server Stopped"})
        } else {
        await StartServer()
          toast({ title: "Server Started"})
        }
        setIsServerRunning(!isServerRunning)
    }
  
    return (
        <div>
      <h1 className="text-3xl font-bold mb-6">Settings</h1>
      <div className="bg-white p-6 rounded-lg shadow-md mb-6">
        <h2 className="text-xl font-semibold mb-4">Streaming Server</h2>
        <div className="space-y-4">
          <div>
            <Label htmlFor="port">Port</Label>
            <Input type="number" id="port" placeholder="8080" className="mt-1" />
          </div>
          <div>
            <Label htmlFor="quality">Default Streaming Quality</Label>
            <select id="quality" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50">
              <option>Original</option>
              <option>1080p</option>
              <option>720p</option>
              <option>480p</option>
            </select>
          </div>
          <div className="flex items-center">
            <input type="checkbox" id="transcoding" className="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50" />
            <Label htmlFor="transcoding" className="ml-2">Enable transcoding</Label>
          </div>
        </div>
        <div className="flex items-center justify-between mt-6">
          <Button onClick={toggleServer}>
            {isServerRunning ? (
              <>
                <Square className="mr-2 h-4 w-4" />
                Stop Server
              </>
            ) : (
              <>
                <Play className="mr-2 h-4 w-4" />
                Start Server
              </>
            )}
          </Button>
          <div className="text-sm">
            Status: <span className={isServerRunning ? "text-green-600" : "text-red-600"}>
              {isServerRunning ? "Running" : "Stopped"}
            </span>
          </div>
        </div>
      </div>
      <div className="bg-white p-6 rounded-lg shadow-md">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Manage Categories</h2>
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Add Category
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Category</DialogTitle>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="name" className="text-right">
                    Name
                  </Label>
                  <Input
                    id="name"
                    value={newCategory}
                    onChange={(e) => setNewCategory(e.target.value)}
                    className="col-span-3"
                  />
                </div>
              </div>
              <div className="flex justify-end">
                <Button onClick={handleAddCategory}>Add Category</Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>
        <ul className="space-y-2">
          {categories.map((category) => (
            <li key={category.ID} className="flex items-center justify-between bg-gray-100 p-3 rounded">
              <span>{category.Name}</span>
              <Button variant="destructive" onClick={() => handleRemoveCategory(category.ID)}>Remove</Button>
            </li>
          ))}
        </ul>
      </div>
    </div>
      )
}