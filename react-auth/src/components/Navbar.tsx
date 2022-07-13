import /*React,*/ { useState } from 'react';
import { Link } from "react-router-dom";
import {
  Navbar,
  NavLink,
  Bars,
  NavMenu,
  //NavBtn,
  NavBtnLink,
  NavLogo
} from './NavbarElements';
//import * as FaIcons from 'react-icons/fa';
import * as AiIcons from 'react-icons/ai';
import { SidebarData } from './SidebarData';
import '../css/Navbar.css';
import { IconContext } from 'react-icons';


const Nav = (props: { name: string, setName: (name: string) => void }) => {


  const [sidebar, setSidebar] = useState(false);

  const showSidebar = () => {
    setSidebar(!sidebar);
    var logo = document.getElementById('nav-logo');

    if(logo) {

      if(sidebar) {
        logo.style.transition = "all 0.2s ease-in-out";
        logo.style.marginLeft = "-150px";
      }
      else {
        logo.style.transition = "all 0.2s ease-in-out";
        logo.style.marginLeft = "0px";
      }
    }

  }

  const logout = async () => {
        await fetch('http://localhost:8000/api/logout', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        });

        props.setName('');
    }

    let menu;

    if (props.name === '') {
        menu = (
            <ul className="navbar-nav me-auto mb-2 mb-md-0">
                <li className="nav-item active">
                    <Link to="/login" className="nav-link">Login</Link>
                </li>
                <li className="nav-item active">
                    <Link to="/register" className="nav-link">Register</Link>
                </li>
            </ul>
        );
    } else {
        menu = (
            <ul className="navbar-nav me-auto mb-2 mb-md-0">
                <li className="nav-item active">
                    <Link to="/login" className="nav-link" onClick={logout}>Logout</Link>
                </li>
            </ul>
        );
    }

    return (
      
      <IconContext.Provider value={{ color: '#000' }}>
          <Navbar>
            <Bars to='#' onClick={showSidebar} className='menu-bars' />
            <NavLogo src={require('../Images/logo.png')} alt='logo' id='nav-logo' />
            <NavMenu>
            <a href='/ejmemo.html' >
                Play
              </a>
              <NavLink to='/about' >
                About
              </NavLink>
              <NavLink to='/contact-us' >
                Contact Us
              </NavLink>
              <NavLink to='/register' >
                Register
              </NavLink>
              <NavBtnLink to='/login' className="nav-login-button">Sign In</NavBtnLink>
            </NavMenu>
        </Navbar>
        <nav className={sidebar ? 'nav-menu active' : 'nav-menu'}>
          <ul className='nav-menu-items' onClick={showSidebar}>
            <li className='navbar-toggle'>
              <Link to='#' className='menu-bars'>
                <AiIcons.AiOutlineClose />
              </Link>
            </li>
            {SidebarData.map((item, index) => {
              return (
                <li key={index} className={item.cName}>
                  <Link to={item.path}>
                    {item.icon}
                    <span>{item.title}</span>
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>
      </IconContext.Provider>
    
    );
};

export default Nav;