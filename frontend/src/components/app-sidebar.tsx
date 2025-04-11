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

/**
 * Represents a menu item on the sidebar.
 */
interface MenuItem {
  /**
   * The parent of menu links. Typically this can be a project team (e.g., TA), or
   * some other parent group ("About").
   */
  project_team: string;
  /**
   * List of the children links. In the case of project teams, this will include
   * libraries; in non teams, can be other ancillary links.
   */
  libraries: string[];
}

/**
 * Renders the app side bar and all its contents.
 * @param {string} setContent - current active state
 * @returns {AppSidebarProps} app sidebar
 */
export function AppSidebar({ setContent, ...props }: AppSidebarProps) {
  const [menuData, setMenuData] = useState<MenuItem[]>([]);

  const extraData: MenuItem = {
    project_team: "About",
    libraries: ["Contribution Guide", "About"],
  };

  useEffect(() => {
    fetch("http://localhost:23456/api/menu/")
      .then((response) => response.text())
      .then((text) => {
        const teamData = JSON.parse(text);
        teamData.push(extraData);
        // add the miscellaneous stuff
        setMenuData(teamData);
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
                defaultOpen={true}
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
                      <SidebarMenuSub className="w-full">
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
export function AppMain({ activeContent }: { activeContent: string }) {
  console.log("Appmain content: ", activeContent);
  return (
    <div className="flex-1 flex flex-col">
      <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
        <div className="flex items-center gap-2 px-4">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className="hidden md:block">
                <BreadcrumbLink href="/">Library of Babel</BreadcrumbLink>
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
