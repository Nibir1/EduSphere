
import { useEffect, useRef, useState } from "react";
import { LogOut } from "lucide-react";
import { useAuth } from "../auth/AuthProvider";

export default function Header() {

  const [open, setOpen] = useState(false);
  const [loggedUser,setLoggedUser] = useState();
  const { logout } = useAuth();
  const panelRef = useRef(null);

  const handleLogout = async () => {
    await logout();   
    setOpen(false);    
  };

  useEffect(() => {
    function handleClick(e) {
      if (panelRef.current && !panelRef.current.contains(e.target)) setOpen(false);
    }
    function handleKey(e) {
      if (e.key === "Escape") setOpen(false);
    }
    if (open) {
      document.addEventListener("mousedown", handleClick);
      document.addEventListener("keydown", handleKey);
    }
    return () => {
      document.removeEventListener("mousedown", handleClick);
      document.removeEventListener("keydown", handleKey);
    };
  }, [open]);

  useEffect(() =>{
    const savedUser = localStorage.getItem("user");
    console.log("savedUser.email",savedUser.email)
    setLoggedUser(JSON.parse(savedUser));
  },[])

  return (
    <header className="border-b border-gray-300 bg-white">
      <div className="mx-auto max-w-7xl px-4 py-2 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">EduSphere</h1>
              <p className="text-sm text-gray-500">Your Personal Academic Advisor</p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <button className="text-sm font-medium text-gray-500 hover:text-gray-900 transition-colors">
              Help
            </button>
            <button
              onClick={() => setOpen((s) => !s)}
              className="h-10 w-10 rounded-full bg-gray-200 text-gray-900 hover:bg-gray-300 transition-colors flex items-center justify-center">
              👤
            </button>
          </div>
        </div>
      </div>

      {open && (
        <div
          ref={panelRef}
          className="absolute right-0 top-14 w-72 bg-white rounded-2xl shadow-xl border border-gray-200 p-5 animate-fadeIn z-50"
        >
          <div className="text-center mb-4">
            {/* Avatar circle */}
            <div className="mx-auto h-16 w-16 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-2xl font-semibold mb-2">
              {loggedUser?.full_name ? loggedUser.full_name.charAt(0).toUpperCase() : "👤"}
            </div>

            {/* User info */}
            <h3 className="text-lg font-semibold text-gray-900">
              {loggedUser?.full_name || "Guest User"}
            </h3>
            <p className="text-sm text-gray-500">{loggedUser?.username }</p>
            <p className="text-sm text-gray-500">{loggedUser?.email || "guest@example.com"}</p>
          </div>

          <div className="border-t border-gray-100 my-4"></div>

          {/* Logout button */}
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center gap-2 py-2 rounded-lg bg-red-500 text-white font-medium hover:bg-red-600 transition"
          >
            <LogOut className="w-4 h-4" />
            Log out
          </button>
        </div>
      )}

    </header>
  )
}
