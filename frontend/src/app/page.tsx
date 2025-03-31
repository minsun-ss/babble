"use client";
import { useState } from "react";
import {
  renderContent,
  renderMain,
  AppSidebar,
} from "@/components/app-sidebar";
import {
  Sidebar,
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";

// Add setContent to props
type AppSidebarProps = React.ComponentProps<typeof Sidebar> & {
  setContent?: (content: string) => void;
};

export default function Page() {
  const [activeContent, setActiveContent] = useState("main");

  return (
    <div className="flex h-screen">
      <SidebarProvider>
        <AppSidebar setContent={setActiveContent}></AppSidebar>
        {renderMain(activeContent)}
      </SidebarProvider>

      {/* Simple custom sidebar */}
      {/* <div className="w-64 bg-gray-100 border-r h-full p-4">
        <div className="py-4">
          <h2 className="text-lg font-semibold mb-4">Navigation</h2>
          <ul className="space-y-2">
            <li>
              <button
                onClick={() => setActiveContent("hello")}
                className={`w-full text-left px-3 py-2 rounded-md transition-colors ${
                  activeContent === "hello"
                    ? "bg-gray-200 font-medium"
                    : "hover:bg-gray-200"
                }`}
              >
                Hello World
              </button>
            </li>
            <li>
              <button
                onClick={() => setActiveContent("main")}
                className={`w-full text-left px-3 py-2 rounded-md transition-colors ${
                  activeContent === "invoices"
                    ? "bg-gray-200 font-medium"
                    : "hover:bg-gray-200"
                }`}
              >
                Invoices
              </button>
            </li>
          </ul>
        </div>
      </div> */}

      {/* Main content area */}
    </div>
  );
}

// export default function Page() {
//   return (
//     <div>
//       <SidebarProvider>
//         <AppSidebar></AppSidebar>
//         <SidebarInset>
//           <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
//             <div className="flex items-center gap-2 px-4">
//               <SidebarTrigger className="-ml-1" />
//               <Separator orientation="vertical" className="mr-2 h-4" />
//               <Breadcrumb>
//                 <BreadcrumbList>
//                   <BreadcrumbItem className="hidden md:block">
//                     <BreadcrumbLink href="#">Babel</BreadcrumbLink>
//                   </BreadcrumbItem>
//                   <BreadcrumbSeparator className="hidden md:block" />
//                   <BreadcrumbItem>
//                     <BreadcrumbPage>Index</BreadcrumbPage>
//                   </BreadcrumbItem>
//                 </BreadcrumbList>
//               </Breadcrumb>{" "}
//             </div>
//           </header>
//           <div className="justify-start items-center gap-2 px-4">
//             <Table>
//               <TableCaption>A list of your recent invoices.</TableCaption>
//               <TableHeader>
//                 <TableRow>
//                   <TableHead className="w-[100px]">Invoice</TableHead>
//                   <TableHead>Status</TableHead>
//                   <TableHead>Method</TableHead>
//                   <TableHead className="text-right">Amount</TableHead>
//                 </TableRow>
//               </TableHeader>
//               <TableBody>
//                 {invoices.map((invoice) => (
//                   <TableRow key={invoice.invoice}>
//                     <TableCell className="font-medium">
//                       {invoice.invoice}
//                     </TableCell>
//                     <TableCell>{invoice.paymentStatus}</TableCell>
//                     <TableCell>{invoice.paymentMethod}</TableCell>
//                     <TableCell className="text-right">
//                       {invoice.totalAmount}
//                     </TableCell>
//                   </TableRow>
//                 ))}
//               </TableBody>
//               <TableFooter>
//                 <TableRow>
//                   <TableCell colSpan={3}>Total</TableCell>
//                   <TableCell className="text-right">$2,500.00</TableCell>
//                 </TableRow>
//               </TableFooter>
//             </Table>

//             <p>Hello world</p>
//             <p>
//               "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
//               eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
//               enim ad minim veniam, quis nostrud exercitation ullamco laboris
//               nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in
//               reprehenderit in voluptate velit esse cillum dolore eu fugiat
//               nulla pariatur. Excepteur sint occaecat cupidatat non proident,
//               sunt in culpa qui officia deserunt mollit anim id est laborum."
//             </p>
//           </div>
//         </SidebarInset>
//       </SidebarProvider>
//     </div>
//   );
// }
