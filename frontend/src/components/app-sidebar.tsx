import { LibraryBig, Minus, Plus } from "lucide-react";
import { renderContent } from "@/components/app-content";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { useState, useEffect } from "react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
} from "@/components/ui/sidebar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";

// Add setContent to AppSidebar's React props
type AppSidebarProps = React.ComponentProps<typeof Sidebar> & {
  setContent?: (content: string) => void;
};

interface MenuItem {
  project_team: string;
  libraries: string[];
}
/**
 * Renders the app side bar.
 * @param {string} setContent - current active state
 * @returns {AppSidebarProps} app sidebar
 */
export function AppSidebar({ setContent, ...props }: AppSidebarProps) {
  const [menuData, setMenuData] = useState<MenuItem[]>([]);

  useEffect(() => {
    fetch("http://localhost:23456/api/menu/")
      .then((response) => response.text())
      .then((text) => {
        setMenuData(JSON.parse(text));
      })
      .catch((error) => {
        console.error("Error: ", error);
      });
  }, []);

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="#">
                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                  <LibraryBig className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-medium">Library of Babel</span>
                  <span className="">v0.2.1</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu>
            {menuData.map((item, index) => (
              <Collapsible
                key={item.project_team}
                defaultOpen={index === 0}
                className="group/collapsible"
              >
                <SidebarMenuItem>
                  <CollapsibleTrigger asChild>
                    <SidebarMenuButton>
                      {item.project_team}{" "}
                      <Plus className="ml-auto group-data-[state=open]/collapsible:hidden" />
                      <Minus className="ml-auto group-data-[state=closed]/collapsible:hidden" />
                    </SidebarMenuButton>
                  </CollapsibleTrigger>
                  {item.libraries?.length ? (
                    <CollapsibleContent>
                      <SidebarMenuSub>
                        {item.libraries.map((library) => (
                          <SidebarMenuSubItem key={library}>
                            <SidebarMenuSubButton asChild>
                              <a
                                href="#"
                                onClick={() =>
                                  setContent && setContent(library)
                                }
                              >
                                {library}
                              </a>
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        ))}
                      </SidebarMenuSub>
                    </CollapsibleContent>
                  ) : null}
                </SidebarMenuItem>
              </Collapsible>
            ))}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  );
}

/**
 * Renders the app main content and breadcrumbs.
 * @param {string} activeContent - current active state
 * @returns {React.ReactElement} app main content and breadcrumb links
 */
export function AppMain(activeContent: string) {
  return (
    <div className="flex-1 flex flex-col">
      <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
        <div className="flex items-center gap-2 px-4">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className="hidden md:block">
                <BreadcrumbLink href="#">Babel</BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator className="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>{activeContent}</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>
      </header>
      <div className="p-6">{renderContent(activeContent)}</div>
    </div>
  );
}
