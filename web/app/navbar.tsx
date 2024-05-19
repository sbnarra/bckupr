'use client';

import {useState} from "react";
import {Navbar as NNavbar, NavbarBrand, NavbarContent, NavbarItem, NavbarMenuToggle, NavbarMenu, NavbarMenuItem, Link, Button} from "@nextui-org/react";

export default function Navbar() {

        const [menuOpen, setMenuOpen] = useState(false)
    return (
    <NNavbar onMenuOpenChange={setMenuOpen}>

        <NavbarContent>
            <NavbarMenuToggle 
                aria-label={menuOpen ? "Close" : "Open"}
                className="sm:hidden"
            />
            <NavbarBrand>
                <p>Bckupr</p>
            </NavbarBrand>
        </NavbarContent>

        <NavbarContent className="hidden sm:flex gab-4" justify="center">
            <NavbarItem isActive>
                <Link aria-current="page" href="/backups">
                    Backups
                </Link>
                <Link color="foreground" href="/settings">
                    Settings
                </Link>
            </NavbarItem>
        </NavbarContent>

        <NavbarContent justify="end">
            <Link color="foreground" href="https://sbnarra.github.io/bckupr/" target="_blank">
                Docs
            </Link>
            <Link color="foreground" href="https://github.com/sbnarra/bckupr" target="_blank">
                Source
            </Link>
        </NavbarContent>

        <NavbarMenu>
            <NavbarMenuItem>
                <Link color="foreground" href="/backups">
                    Backups
                </Link>
            </NavbarMenuItem>
            <NavbarMenuItem>
                <Link color="danger" href="/settings">
                    Settings
                </Link>
            </NavbarMenuItem>
        </NavbarMenu>

    </NNavbar>)
}