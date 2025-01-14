﻿//------------------------------------------------------------------------------
// <auto-generated>
//     This code was generated by a tool.
//     Runtime Version:4.0.30319.0
//
//     Changes to this file may cause incorrect behavior and will be lost if
//     the code is regenerated.
// </auto-generated>
//------------------------------------------------------------------------------

namespace MageServer.Resources {
    using System;
    
    
    /// <summary>
    ///   A strongly-typed resource class, for looking up localized strings, etc.
    /// </summary>
    // This class was auto-generated by the StronglyTypedResourceBuilder
    // class via a tool like ResGen or Visual Studio.
    // To add or remove a member, edit your .ResX file then rerun ResGen
    // with the /str option, or rebuild your VS project.
    [global::System.CodeDom.Compiler.GeneratedCodeAttribute("System.Resources.Tools.StronglyTypedResourceBuilder", "4.0.0.0")]
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute()]
    [global::System.Runtime.CompilerServices.CompilerGeneratedAttribute()]
    public class Strings_MySQL {
        
        private static global::System.Resources.ResourceManager resourceMan;
        
        private static global::System.Globalization.CultureInfo resourceCulture;
        
        [global::System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("Microsoft.Performance", "CA1811:AvoidUncalledPrivateCode")]
        internal Strings_MySQL() {
        }
        
        /// <summary>
        ///   Returns the cached ResourceManager instance used by this class.
        /// </summary>
        [global::System.ComponentModel.EditorBrowsableAttribute(global::System.ComponentModel.EditorBrowsableState.Advanced)]
        public static global::System.Resources.ResourceManager ResourceManager {
            get {
                if (object.ReferenceEquals(resourceMan, null)) {
                    global::System.Resources.ResourceManager temp = new global::System.Resources.ResourceManager("MageServer.Resources.Strings_MySQL", typeof(Strings_MySQL).Assembly);
                    resourceMan = temp;
                }
                return resourceMan;
            }
        }
        
        /// <summary>
        ///   Overrides the current thread's CurrentUICulture property for all
        ///   resource lookups using this strongly typed resource class.
        /// </summary>
        [global::System.ComponentModel.EditorBrowsableAttribute(global::System.ComponentModel.EditorBrowsableState.Advanced)]
        public static global::System.Globalization.CultureInfo Culture {
            get {
                return resourceCulture;
            }
            set {
                resourceCulture = value;
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to Database={0}; Data Source={1}; Port={2}; User Id={3}; Password={4};.
        /// </summary>
        public static string ConnectionString {
            get {
                return ResourceManager.GetString("ConnectionString", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to Error connecting to the MySQL server. You should restart the game and MySQL server..
        /// </summary>
        public static string Error_Connecting {
            get {
                return ResourceManager.GetString("Error_Connecting", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to DELETE FROM `characters` WHERE accountid = @accountid AND name = @name.
        /// </summary>
        public static string NonQuery_Delete_Character_Delete {
            get {
                return ResourceManager.GetString("NonQuery_Delete_Character_Delete", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to DELETE FROM `online_characters` WHERE charid = @charid.
        /// </summary>
        public static string NonQuery_Delete_Character_SetOffline {
            get {
                return ResourceManager.GetString("NonQuery_Delete_Character_SetOffline", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to DELETE FROM  `online_accounts`.
        /// </summary>
        public static string NonQuery_Delete_OnlineAccounts_SetAllOffline {
            get {
                return ResourceManager.GetString("NonQuery_Delete_OnlineAccounts_SetAllOffline", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to DELETE FROM `online_accounts` WHERE accountid = @accountid.
        /// </summary>
        public static string NonQuery_Delete_OnlineAccounts_SetOffline {
            get {
                return ResourceManager.GetString("NonQuery_Delete_OnlineAccounts_SetOffline", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to DELETE FROM `online_characters`.
        /// </summary>
        public static string NonQuery_Delete_OnlineCharacters_SetAllOffline {
            get {
                return ResourceManager.GetString("NonQuery_Delete_OnlineCharacters_SetAllOffline", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to DELETE FROM `character_statistics` WHERE charid = @charid.
        /// </summary>
        public static string NonQuery_Delete_StatisticsOverall_DeleteByCharId {
            get {
                return ResourceManager.GetString("NonQuery_Delete_StatisticsOverall_DeleteByCharId", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to DELETE FROM `character_statistics_weekly` WHERE charid = @charid.
        /// </summary>
        public static string NonQuery_Delete_StatisticsWeekly_DeleteByCharId {
            get {
                return ResourceManager.GetString("NonQuery_Delete_StatisticsWeekly_DeleteByCharId", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to INSERT INTO `characters` (accountid, slot, name, agility, constitution, memory, reasoning, discipline, empathy, intuition, presence, quickness, strength, spent_stat, bonus_stat, bonus_spent, list_1, list_2, list_3, list_4, list_5, list_6, list_7, list_8, list_9, list_10, list_level_1, list_level_2, list_level_3, list_level_4, list_level_5, list_level_6, list_level_7, list_level_8, list_level_9, list_level_10, class, level, spell_picks, model, experience, spell_key_1, spell_key_2, spell_key_3, spell_key_4, s [rest of string was truncated]&quot;;.
        /// </summary>
        public static string NonQuery_Insert_Character_SaveNew {
            get {
                return ResourceManager.GetString("NonQuery_Insert_Character_SaveNew", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to INSERT INTO `matches` (arenaid, tableid, creation_time, player_count, highest_player_count, max_players, current_state, end_state, short_name, long_name, founder_charid, duration, level_range, mode, rules) VALUES (@arenaid, @tableid, @creation_time, @player_count, @highest_player_count, @max_players, @current_state, @end_state, @short_name, @long_name, @founder_charid, @duration, @level_range, @mode, @rules).
        /// </summary>
        public static string NonQuery_Insert_Matches_New {
            get {
                return ResourceManager.GetString("NonQuery_Insert_Matches_New", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to INSERT INTO `online_accounts` (accountid) VALUES (@accountid).
        /// </summary>
        public static string NonQuery_Insert_OnlineAccounts_SetOnline {
            get {
                return ResourceManager.GetString("NonQuery_Insert_OnlineAccounts_SetOnline", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to INSERT INTO `online_characters` (charid, tableid, arenaid, arenashortname) VALUES (@charid, @tableid, @arenaid, @arenashortname) ON DUPLICATE KEY UPDATE tableid=@tableid, arenaid=@arenaid, arenashortname=@arenashortname.
        /// </summary>
        public static string NonQuery_InsertUpdate_OnlineCharacters_SetOnline {
            get {
                return ResourceManager.GetString("NonQuery_InsertUpdate_OnlineCharacters_SetOnline", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to INSERT INTO `character_statistics` (charid, hidden, kills, deaths, raises, damagedone, damagetaken, healingdone, healingtaken, wins, losses) VALUES (@charid, @hidden, @kills, @deaths, @raises, @damagedone, @damagetaken, @healingdone, @healingtaken, @wins, @losses) ON DUPLICATE KEY UPDATE hidden=@hidden, kills=kills+@kills, deaths=deaths+@deaths, raises=raises+@raises, damagedone=damagedone+@damagedone, damagetaken=damagetaken+@damagetaken, healingdone=healingdone+@healingdone, healingtaken=healingtaken+@hea [rest of string was truncated]&quot;;.
        /// </summary>
        public static string NonQuery_InsertUpdate_StatisticsOverall_Update {
            get {
                return ResourceManager.GetString("NonQuery_InsertUpdate_StatisticsOverall_Update", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to INSERT INTO `character_statistics_weekly` (charid, date, hidden, kills, deaths, raises, damagedone, damagetaken, healingdone, healingtaken, wins, losses) VALUES (@charid, @date, @hidden, @kills, @deaths, @raises, @damagedone, @damagetaken, @healingdone, @healingtaken, @wins, @losses) ON DUPLICATE KEY UPDATE hidden=@hidden, kills=kills+@kills, deaths=deaths+@deaths, raises=raises+@raises, damagedone=@damagedone+@damagedone, damagetaken=damagetaken+@damagetaken, healingdone=healingdone+@healingdone, healingta [rest of string was truncated]&quot;;.
        /// </summary>
        public static string NonQuery_InsertUpdate_StatisticsWeekly_Update {
            get {
                return ResourceManager.GetString("NonQuery_InsertUpdate_StatisticsWeekly_Update", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to UPDATE `characters` SET name=@name, agility=@agility, constitution=@constitution, memory=@memory, reasoning=@reasoning, discipline=@discipline, empathy=@empathy, intuition=@intuition, presence=@presence, quickness=@quickness, strength=@strength, spent_stat=@spent_stat, bonus_stat=@bonus_stat, bonus_spent=@bonus_spent, list_1=@list_1, list_2=@list_2, list_3=@list_3, list_4=@list_4, list_5=@list_5, list_6=@list_6, list_7=@list_7, list_8=@list_8, list_9=@list_9, list_10=@list_10, list_level_1=@list_level_1, li [rest of string was truncated]&quot;;.
        /// </summary>
        public static string NonQuery_Update_Character_SaveExisting {
            get {
                return ResourceManager.GetString("NonQuery_Update_Character_SaveExisting", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to UPDATE `server_settings` SET exp_multiplier = @exp_multiplier.
        /// </summary>
        public static string NonQuery_Update_ServerSettings_SetExpMultiplier {
            get {
                return ResourceManager.GetString("NonQuery_Update_ServerSettings_SetExpMultiplier", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to SELECT `serial` FROM `banned_serials` WHERE serial = @serial LIMIT 1.
        /// </summary>
        public static string Query_Select_Account_IsSerialBanned {
            get {
                return ResourceManager.GetString("Query_Select_Account_IsSerialBanned", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to SELECT * FROM `characters` WHERE accountid = @accountid AND slot = @slot LIMIT 1.
        /// </summary>
        public static string Query_Select_Character_FindByAccountIdAndSlot {
            get {
                return ResourceManager.GetString("Query_Select_Character_FindByAccountIdAndSlot", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to SELECT * FROM `characters` WHERE name = @name LIMIT 1.
        /// </summary>
        public static string Query_Select_Character_FindByName {
            get {
                return ResourceManager.GetString("Query_Select_Character_FindByName", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to SELECT * FROM `characters` WHERE name = @name AND accountid = @accountid LIMIT 1.
        /// </summary>
        public static string Query_Select_Character_FindByNameAndAccountId {
            get {
                return ResourceManager.GetString("Query_Select_Character_FindByNameAndAccountId", resourceCulture);
            }
        }
        
        /// <summary>
        ///   Looks up a localized string similar to SELECT * FROM `characters` WHERE class = @class AND oplevel = 0 ORDER BY `experience` DESC LIMIT 10.
        /// </summary>
        public static string Query_Select_Study_GetHighScoreList {
            get {
                return ResourceManager.GetString("Query_Select_Study_GetHighScoreList", resourceCulture);
            }
        }
    }
}
