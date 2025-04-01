"use client";
import { useState } from "react";
import { AppMain, AppSidebar } from "@/components/app-sidebar";
import { SidebarProvider } from "@/components/ui/sidebar";

export default function Page() {
  const [activeContent, setActiveContent] = useState("index");

  return (
    <div className="flex h-screen">
      <SidebarProvider>
        <AppSidebar setContent={setActiveContent}></AppSidebar>
        {AppMain(activeContent)}
      </SidebarProvider>
    </div>
  );
}
