using System.Text.RegularExpressions;

namespace DevSmtp.Core.Models
{
    public sealed class Email : IEquatable<Email>
    {
        private static readonly Regex EmailRegex =
            new("^[\\w!#$%&’*+/=?`{|}~^-]+(?:\\.[\\w!#$%&’*+/=?`{|}~^-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6}$",
                RegexOptions.Compiled | RegexOptions.IgnoreCase);

        public Email(string value)
        {
            this.Value = value;
            this.Validate();
        }

        public string Value { get; }
        
        public static Email From(string value) => new(value);

        private void Validate()
        {
            if(!EmailRegex.IsMatch(this.Value))
            {
                throw new FormatException($"{this.Value} is not a valid email address.");
            }
        }

        public bool Equals(Email? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return this.Value == other.Value;
        }

        public override bool Equals(object? obj)
        {
            return ReferenceEquals(this, obj) || obj is Email other && Equals(other);
        }

        public override int GetHashCode()
        {
            return this.Value.GetHashCode();
        }
    }
}
