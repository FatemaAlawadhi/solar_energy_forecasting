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
        <button
          aria-controls="sidebar"
          onClick={(e) => {
            e.stopPropagation();
            props.setSidebarOpen(!props.sidebarOpen);
          }}
          className="block lg:hidden"
        >
          <svg
            className="fill-current"
            width="20"
            height="18"
            viewBox="0 0 20 18"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M19 8.175H2.98748L9.36248 1.6875C9.69998 1.35 9.69998 0.825 9.36248 0.4875C9.02498 0.15 8.49998 0.15 8.16248 0.4875L0.399976 8.3625C0.0624756 8.7 0.0624756 9.225 0.399976 9.5625L8.16248 17.4375C8.31248 17.5875 8.53748 17.7 8.76248 17.7C8.98748 17.7 9.21248 17.625 9.36248 17.4375C9.69998 17.1 9.69998 16.575 9.36248 16.2375L2.98748 9.75H19C19.45 9.75 19.825 9.375 19.825 8.925C19.825 8.475 19.45 8.175 19 8.175Z"
              fill=""
            />
          </svg>
        </button>

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
