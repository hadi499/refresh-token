import axios from "axios";
import jwt_decode from "jwt-decode";
import dayjs from "dayjs";
import { useContext } from "react";
import AuthContext from "../context/AuthContext";

const baseURL = "http://localhost:5000";

const useAxios = () => {
  const { authTokens} = useContext(AuthContext);

  const axiosInstance = axios.create({
    baseURL,
    headers: { Authorization: `Bearer ${authTokens?.accessToken}` },
  });

  axiosInstance.interceptors.request.use(async (req) => {
    const user = jwt_decode(authTokens.accessToken);
    const isExpired = dayjs.unix(user.exp).diff(dayjs()) < 1;

    if (!isExpired) return req;
    
    //${baseURL}/api/users/refresh/ berjalan di nodejs tapi tidak di go
    // dibawah ini untuk go
    const response = await axios.post(`${baseURL}/api/users/refresh`, {
      token: authTokens.refreshToken,
    });

    localStorage.setItem("authTokens", JSON.stringify(response.data));   

    req.headers.Authorization = `Bearer ${response.data.accessToken}`;
    return req;
  });

  return axiosInstance;
};

export default useAxios;
