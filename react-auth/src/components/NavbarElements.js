import { FaBars } from 'react-icons/fa';
import { NavLink as Link } from 'react-router-dom';
import styled from 'styled-components';

export const Navbar = styled.nav`
  background: #fff;
  height: 65px;
  display: flex;
  justify-content: space-between;
  padding: 0.5rem calc((100vw - 1000px) / 2);
  z-index: 10;

  /* Third Nav */
  /* justify-content: flex-start; */
`;

export const NavLink = styled(Link)`
  color: #999;
  display: flex;
  align-items: center;
  text-decoration: none;
  padding: 0 1rem;
  height: 100%;
  cursor: pointer;

  &.active {
    color: #ff4b46;
  }
`;

export const Bars = styled(FaBars)`

  color: #fff;
  display: block;
  position: absolute;
  text-align: center;
  top: 0;
  left: 50px;
  transform: translate(-100%, 60%);
  font-size: 1.8rem;
  cursor: pointer;
`;

export const NavMenu = styled.div`
  display: flex;
  align-items: center;
  margin-right: -24px;

  /* Second Nav */
  /* margin-right: 24px; */


  /* width: 100vw;
  white-space: nowrap; */

  @media screen and (max-width: 1200px) {
    display: none;
  }
`;

export const NavBtn = styled.nav`
  display: flex;
  align-items: center;
  margin-right: 24px;


  justify-content: flex-end;
  width: 100vw; */

  @media screen and (max-width: 1200px) {
    display: none;
  }
`;

export const NavBtnLink = styled(Link)`
  border-radius: 4px;
  background: #256ce1;
  padding: 10px 22px;
  color: #fff;
  outline: none;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  text-decoration: none;
  background-color: #ff4b46;

  margin-left: 24px;

  &:hover {
    transition: all 0.2s ease-in-out;
    background: #ff4b46;
    color: #fff;
  }

`;

export const NavLogo = styled.img` {
  display: block;
  position: flex;
  height: 95%;
  margin-left: -150px;
  @media screen and (max-width: 1200px) {
    display: none;
  }
}`;

