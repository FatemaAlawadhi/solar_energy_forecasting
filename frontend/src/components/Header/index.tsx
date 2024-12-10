import Link from "next/link";
import DarkModeSwitcher from "./DarkModeSwitcher";
import Image from "next/image";

const Header = (props: {
  sidebarOpen: string | boolean | undefined;
  setSidebarOpen: (arg0: boolean) => void;
}) => {
  return (
    <header className="sticky top-0 z-999 flex w-full border-b border-stroke bg-white dark:border-stroke-dark dark:bg-gray-dark">
      <div className="flex flex-grow items-center justify-between px-4 py-5 shadow-2 md:px-5 2xl:px-10">
        
        <div className="flex items-center gap-3">
          <Link href="/" className="flex items-center gap-3">
            <div className="w-[50px] h-[50px] rounded-full overflow-hidden">
              <Image
                width={40}
                height={40}
                src={"/images/logo/bapco-logo.png"}
                alt="Logo"
                priority
                style={{ width: "100%", height: "100%", objectFit: "cover" }}
              />
            </div>
            <h3 className="mb-0.5 text-heading-5 font-bold text-dark dark:text-white"> 
              Bapco Energies
            </h3>
          </Link>
        </div>

        <DarkModeSwitcher />
        
      </div>
    </header>
  );
};

export default Header;
