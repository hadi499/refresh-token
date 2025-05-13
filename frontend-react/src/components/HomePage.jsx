import  { useEffect, useState } from "react";
import useAxios from "../utils/useAxios";

const HomePage = () => {

  const [users, setUsers] = useState([]);
  let api = useAxios();

  useEffect(() => {
    getUsers();
  }, []);

  let getUsers = async () => {
    let response = await api.get("/api/users");

    if (response.status === 200) {
      setUsers(response.data);
    }
  };
  

  return (
    <div className="flex justify-center h-screen mt-12 ">
      <div className="w-1/3 ">
        <div className="flex flex-col gap-3">
          {users.map((user) => (
            // gunakan user._id untuk nodejs mongodb , user.id untuk golang
            <div key={user._id}>
              <h1 className="text-2xl font-bold">{user.name}</h1>
              <p className="text-lg font-semibold">{user.email}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default HomePage;
